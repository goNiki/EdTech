package quiz

import (
	"context"
	"errors"
	"fmt"
	"time"

	"edtech/internal/domain"
	"edtech/internal/infrastructure/txmanager"
	errorsAPP "edtech/pkg/errors"

	"github.com/jackc/pgx/v5"
)

func (r *repositoryImpl) GetStudentHomeworkFeedback(ctx context.Context, userID, lessonID int64) (*domain.StudentHomeworkFeedback, error) {
	const op = "repository.quiz.GetStudentHomeworkFeedback"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	// 1. Поиск последней попытки студента по уроку (с приоритетом essay-квиза, если таковой есть)
	attemptQuery := `
		SELECT qa.id, qa.started_at, qa.completed_at
		FROM quiz_attempts qa
		JOIN quizzes qz ON qz.id = qa.quiz_id
		WHERE qz.lesson_id = $1 
		  AND qa.user_id = $2
		  AND qz.deleted_at IS NULL
		ORDER BY CASE WHEN qz.type = 'essay' THEN 1 ELSE 2 END ASC, qa.id DESC
		LIMIT 1
	`

	var (
		attemptID   int64
		startedAt   time.Time
		completedAt *time.Time
	)

	err := q.QueryRow(ctx, attemptQuery, lessonID, userID).Scan(&attemptID, &startedAt, &completedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &domain.StudentHomeworkFeedback{
				HasSubmission: false,
				Status:        "not_submitted",
				Answers:       make([]domain.HomeworkAnswerDetail, 0),
			}, nil
		}
		return nil, fmt.Errorf("%s: query attempt: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	submittedAt := startedAt
	if completedAt != nil {
		submittedAt = *completedAt
	}

	// 2. Получение ответов студента, баллов, рецензий и данных проверившего преподавателя
	answersQuery := `
		SELECT 
			qaa.id AS answer_id,
			COALESCE(qq.question_text, '') AS question_text,
			COALESCE(qaa.user_answer, '') AS student_answer,
			COALESCE(qaa.points, 0) AS points,
			COALESCE(NULLIF(qz.points, 0), 25) AS max_points,
			qaa.is_correct,
			qaa.feedback,
			qaa.graded_at,
			COALESCE(tu.id, cu.id) AS teacher_id,
			COALESCE(NULLIF(TRIM(tu.first_name || ' ' || tu.last_name), ''), tu.username, NULLIF(TRIM(cu.first_name || ' ' || cu.last_name), ''), cu.username) AS teacher_name,
			COALESCE(tu.avatar_url, cu.avatar_url) AS teacher_avatar_url
		FROM quiz_attempt_answers qaa
		JOIN quiz_questions qq ON qq.id = qaa.question_id
		JOIN quiz_attempts qa ON qa.id = qaa.attempt_id
		JOIN quizzes qz ON qz.id = qa.quiz_id
		LEFT JOIN users tu ON tu.id = qaa.graded_by
		LEFT JOIN courses c ON c.id = qz.course_id
		LEFT JOIN users cu ON cu.id = c.created_by
		WHERE qaa.attempt_id = $1
		ORDER BY qq.position ASC, qaa.id ASC
	`

	rows, err := q.Query(ctx, answersQuery, attemptID)
	if err != nil {
		return nil, fmt.Errorf("%s: query answers: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}
	defer rows.Close()

	answers := make([]domain.HomeworkAnswerDetail, 0)
	var latestGradedAt *time.Time
	var teacher *domain.HomeworkTeacherInfo
	hasUngraded := false

	for rows.Next() {
		var (
			ans           domain.HomeworkAnswerDetail
			gradedAt      *time.Time
			teacherID     *int64
			teacherName   *string
			teacherAvatar *string
		)

		if err := rows.Scan(
			&ans.AnswerID,
			&ans.QuestionText,
			&ans.StudentAnswer,
			&ans.Points,
			&ans.MaxPoints,
			&ans.IsCorrect,
			&ans.Feedback,
			&gradedAt,
			&teacherID,
			&teacherName,
			&teacherAvatar,
		); err != nil {
			return nil, fmt.Errorf("%s: scan answer: %w: %w", op, errorsAPP.ErrInternalDB, err)
		}

		if ans.IsCorrect == nil {
			hasUngraded = true
		} else if gradedAt != nil {
			if latestGradedAt == nil || gradedAt.After(*latestGradedAt) {
				latestGradedAt = gradedAt
			}
		}

		if teacher == nil && teacherID != nil && *teacherID > 0 && teacherName != nil && *teacherName != "" {
			teacher = &domain.HomeworkTeacherInfo{
				ID:        *teacherID,
				Name:      *teacherName,
				AvatarURL: teacherAvatar,
			}
		}

		answers = append(answers, ans)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows err: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	if len(answers) == 0 {
		return &domain.StudentHomeworkFeedback{
			HasSubmission: false,
			Status:        "not_submitted",
			Answers:       make([]domain.HomeworkAnswerDetail, 0),
		}, nil
	}

	status := "graded"
	if hasUngraded {
		status = "pending"
		latestGradedAt = nil
		teacher = nil
	} else if latestGradedAt == nil && completedAt != nil {
		latestGradedAt = completedAt
	}

	return &domain.StudentHomeworkFeedback{
		HasSubmission: true,
		Status:        status,
		AttemptID:     &attemptID,
		SubmittedAt:   &submittedAt,
		GradedAt:      latestGradedAt,
		Teacher:       teacher,
		Answers:       answers,
	}, nil
}
