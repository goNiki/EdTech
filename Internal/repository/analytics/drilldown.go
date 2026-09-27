package analytics

import (
	"context"
	"fmt"
	"time"

	"edtech/internal/domain"
	"edtech/internal/infrastructure/db"
	errorsAPP "edtech/pkg/errors"
)

func (r *repositoryAnalytics) GetStudentDrilldown(ctx context.Context, q db.QueryExecutor, userID, courseID int64) (*domain.StudentDrilldownReport, error) {
	const op = "repository.analytics.GetStudentDrilldown"

	// 1. Get student info
	studentQuery := `
		SELECT id, email, username, COALESCE(first_name, ''), COALESCE(last_name, ''), avatar_url, bio, COALESCE(role::text, 'student'), email_verified, created_at, last_login_at
		FROM users WHERE id = $1 AND deleted_at IS NULL
	`
	var student domain.User
	var roleStr string
	err := q.QueryRow(ctx, studentQuery, userID).Scan(
		&student.ID, &student.Email, &student.Username,
		&student.FirstName, &student.LastName, &student.AvatarURL,
		&student.Bio, &roleStr, &student.EmailVerified,
		&student.CreatedAt, &student.LastLoginAt,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: get user: %w: %w", op, errorsAPP.ErrUserNotFound, err)
	}
	student.Role = domain.Role(roleStr)

	// 2. Get lesson progress logs
	lessonLogsQuery := `
		SELECT 
			l.id, l.title, l.type, 
			COALESCE(lp.status, 'not_started'), 
			lp.score, lp.completed_at, lp.last_accessed_at
		FROM lessons l
		LEFT JOIN lesson_progress lp ON l.id = lp.lesson_id AND lp.user_id = $1
		WHERE l.course_id = $2 AND l.deleted_at IS NULL
		ORDER BY l.position ASC
	`
	rows, err := q.Query(ctx, lessonLogsQuery, userID, courseID)
	if err != nil {
		return nil, fmt.Errorf("%s: lesson logs: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}
	defer rows.Close()

	var lessonLogs []domain.StudentLessonLog
	completedCount := 0
	totalScoreSum := 0
	scoreCount := 0

	for rows.Next() {
		var log domain.StudentLessonLog
		if err := rows.Scan(
			&log.LessonID, &log.LessonTitle, &log.LessonType,
			&log.Status, &log.Score, &log.CompletedAt, &log.LastAccessedAt,
		); err != nil {
			return nil, fmt.Errorf("%s: scan lesson log: %w: %w", op, errorsAPP.ErrInternalDB, err)
		}
		if log.Status == "completed" {
			completedCount++
		}
		if log.Score != nil && *log.Score > 0 {
			totalScoreSum += *log.Score
			scoreCount++
		}
		lessonLogs = append(lessonLogs, log)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows err: %w", op, err)
	}

	// 3. Get all quiz attempts with answers
	attemptsQuery := `
		SELECT 
			qa.id AS attempt_id,
			qz.id AS quiz_id,
			qz.title AS quiz_title,
			ROW_NUMBER() OVER (PARTITION BY qa.quiz_id ORDER BY qa.started_at ASC) AS attempt_num,
			qa.score::int,
			qa.passed,
			qa.started_at,
			qa.completed_at,
			qq.id AS question_id,
			qq.question_text AS question_text,
			COALESCE(qaa.user_answer, ans.text, '') AS chosen_answer,
			qaa.is_correct,
			COALESCE(qaa.points, 0) AS points,
			COALESCE(qz.points, 10) AS max_points,
			qaa.feedback,
			qaa.attachment_url,
			COALESCE(ans.explain, '') AS explanation
		FROM quiz_attempts qa
		JOIN quizzes qz ON qa.quiz_id = qz.id
		JOIN quiz_questions qq ON qz.id = qq.quiz_id
		LEFT JOIN quiz_attempt_answers qaa ON qa.id = qaa.attempt_id AND qq.id = qaa.question_id
		LEFT JOIN quiz_answers ans ON qaa.answer_id = ans.id
		WHERE qz.course_id = $1 AND qa.user_id = $2
		ORDER BY qa.started_at ASC, qq.position ASC
	`

	attRows, err := q.Query(ctx, attemptsQuery, courseID, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: attempts query: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}
	defer attRows.Close()

	attemptsMap := make(map[int64]*domain.StudentTestAttemptLog)
	var attemptsOrder []int64

	for attRows.Next() {
		var attemptID, quizID, questionID int64
		var quizTitle, questionText, chosenAnswer, explanation string
		var attemptNum, score, points, maxPoints int
		var passed bool
		var isCorrect *bool
		var feedback, attachmentURL *string
		var startedAt time.Time
		var completedAt *time.Time

		if err := attRows.Scan(
			&attemptID, &quizID, &quizTitle, &attemptNum, &score, &passed,
			&startedAt, &completedAt, &questionID, &questionText, &chosenAnswer,
			&isCorrect, &points, &maxPoints, &feedback, &attachmentURL, &explanation,
		); err != nil {
			return nil, fmt.Errorf("%s: scan attempt row: %w: %w", op, errorsAPP.ErrInternalDB, err)
		}

		attLog, exists := attemptsMap[attemptID]
		if !exists {
			attLog = &domain.StudentTestAttemptLog{
				AttemptID:     attemptID,
				QuizID:        quizID,
				QuizTitle:     quizTitle,
				AttemptNumber: attemptNum,
				Score:         score,
				Passed:        passed,
				StartedAt:     startedAt,
				CompletedAt:   completedAt,
				Answers:       []domain.StudentAttemptAnswerDetail{},
			}
			attemptsMap[attemptID] = attLog
			attemptsOrder = append(attemptsOrder, attemptID)
		}

		attLog.Answers = append(attLog.Answers, domain.StudentAttemptAnswerDetail{
			QuestionID:    questionID,
			QuestionText:  questionText,
			ChosenAnswer:  chosenAnswer,
			IsCorrect:     isCorrect,
			Points:        points,
			MaxPoints:     maxPoints,
			Feedback:      feedback,
			AttachmentURL: attachmentURL,
			Explanation:   explanation,
		})
	}

	if err := attRows.Err(); err != nil {
		return nil, fmt.Errorf("%s: attRows err: %w", op, err)
	}

	var testAttempts []domain.StudentTestAttemptLog
	for _, id := range attemptsOrder {
		if a, ok := attemptsMap[id]; ok {
			testAttempts = append(testAttempts, *a)
		}
	}

	overallProgress := 0.0
	if len(lessonLogs) > 0 {
		overallProgress = float64(completedCount*100) / float64(len(lessonLogs))
	}
	avgScore := 100.0
	if scoreCount > 0 {
		avgScore = float64(totalScoreSum) / float64(scoreCount)
	}

	return &domain.StudentDrilldownReport{
		Student:         student,
		CourseID:        courseID,
		OverallProgress: overallProgress,
		AvgScore:        avgScore,
		LessonLogs:      lessonLogs,
		TestAttempts:    testAttempts,
	}, nil
}
