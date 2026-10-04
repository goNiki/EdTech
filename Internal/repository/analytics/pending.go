package analytics

import (
	"context"
	"fmt"
	"time"

	"edtech/internal/domain"
	"edtech/internal/infrastructure/txmanager"
	errorsAPP "edtech/pkg/errors"
)

func (r *repositoryAnalytics) ListPendingHomeworks(ctx context.Context, courseID int64, limit, offset int64) ([]domain.PendingHomeworkItem, int64, error) {
	const op = "repository.analytics.ListPendingHomeworks"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	countQuery := `
		SELECT COUNT(*)
		FROM quiz_attempt_answers qaa
		JOIN quiz_attempts qa ON qaa.attempt_id = qa.id
		JOIN quizzes qz ON qa.quiz_id = qz.id
		WHERE qz.course_id = $1 AND qaa.is_correct IS NULL
	`
	var total int64
	err := q.QueryRow(ctx, countQuery, courseID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: count: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	query := `
		SELECT 
			qa.id AS attempt_id,
			qaa.id AS answer_id,
			u.id AS student_id,
			COALESCE(NULLIF(TRIM(u.first_name || ' ' || u.last_name), ''), u.username) AS student_name,
			u.email AS student_email,
			u.username AS student_username,
			c.id AS course_id,
			c.title AS course_title,
			l.id AS lesson_id,
			l.title AS lesson_title,
			qq.id AS question_id,
			qq.question_text AS question_text,
			COALESCE(qaa.user_answer, '') AS student_answer,
			qaa.attachment_url,
			COALESCE(qz.points, 10) AS max_points,
			qa.created_at AS submitted_at
		FROM quiz_attempt_answers qaa
		JOIN quiz_attempts qa ON qaa.attempt_id = qa.id
		JOIN quizzes qz ON qa.quiz_id = qz.id
		JOIN courses c ON qz.course_id = c.id
		JOIN lessons l ON qz.lesson_id = l.id
		JOIN quiz_questions qq ON qaa.question_id = qq.id
		JOIN users u ON qa.user_id = u.id
		WHERE qz.course_id = $1 AND qaa.is_correct IS NULL
		ORDER BY qa.created_at ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := q.Query(ctx, query, courseID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: query: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}
	defer rows.Close()

	var items []domain.PendingHomeworkItem
	for rows.Next() {
		var item domain.PendingHomeworkItem
		var submittedAt time.Time

		if err := rows.Scan(
			&item.AttemptID,
			&item.AnswerID,
			&item.StudentID,
			&item.StudentName,
			&item.StudentEmail,
			&item.StudentUsername,
			&item.CourseID,
			&item.CourseTitle,
			&item.LessonID,
			&item.LessonTitle,
			&item.QuestionID,
			&item.QuestionText,
			&item.StudentAnswer,
			&item.AttachmentURL,
			&item.MaxPoints,
			&submittedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("%s: scan: %w: %w", op, errorsAPP.ErrInternalDB, err)
		}

		item.SubmittedAt = submittedAt
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("%s: rows error: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	if items == nil {
		items = []domain.PendingHomeworkItem{}
	}

	return items, total, nil
}
