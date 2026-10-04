package analytics

import (
	"context"
	"fmt"
	"time"

	"edtech/internal/domain"
	"edtech/internal/infrastructure/txmanager"
	errorsAPP "edtech/pkg/errors"
)

func (r *repositoryAnalytics) ListTeacherPendingHomeworks(
	ctx context.Context,
	teacherID int64,
	courseID int64,
	limit, offset int64,
) ([]domain.PendingHomeworkItem, int64, []domain.CoursePendingSummaryItem, error) {
	const op = "repository.analytics.ListTeacherPendingHomeworks"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	// 1. Fetch courses summary with pending counts for this teacher
	summaryQuery := `
		SELECT 
			c.id AS course_id, 
			c.title AS course_title, 
			COUNT(qaa.id) FILTER (WHERE qaa.is_correct IS NULL) AS pending_count
		FROM courses c
		LEFT JOIN quizzes qz ON qz.course_id = c.id
		LEFT JOIN quiz_attempts qa ON qa.quiz_id = qz.id
		LEFT JOIN quiz_attempt_answers qaa ON qaa.attempt_id = qa.id AND qaa.is_correct IS NULL
		WHERE (c.created_by = $1 OR EXISTS (
			SELECT 1 FROM users_courses uc 
			WHERE uc.course_id = c.id AND uc.user_id = $1 AND uc.role IN ('teacher', 'creator')
		))
		  AND c.deleted_at IS NULL
		GROUP BY c.id, c.title
		ORDER BY pending_count DESC, c.title ASC
	`

	summaryRows, err := q.Query(ctx, summaryQuery, teacherID)
	if err != nil {
		return nil, 0, nil, fmt.Errorf("%s: summary query: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}
	defer summaryRows.Close()

	summary := make([]domain.CoursePendingSummaryItem, 0)
	for summaryRows.Next() {
		var item domain.CoursePendingSummaryItem
		if err := summaryRows.Scan(&item.CourseID, &item.CourseTitle, &item.PendingCount); err != nil {
			return nil, 0, nil, fmt.Errorf("%s: summary scan: %w: %w", op, errorsAPP.ErrInternalDB, err)
		}
		summary = append(summary, item)
	}
	if err := summaryRows.Err(); err != nil {
		return nil, 0, nil, fmt.Errorf("%s: summary rows: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	// 2. Count total pending items matching filters
	countQuery := `
		SELECT COUNT(*)
		FROM quiz_attempt_answers qaa
		JOIN quiz_attempts qa ON qaa.attempt_id = qa.id
		JOIN quizzes qz ON qa.quiz_id = qz.id
		JOIN courses c ON qz.course_id = c.id
		WHERE qaa.is_correct IS NULL
		  AND (c.created_by = $1 OR EXISTS (
			  SELECT 1 FROM users_courses uc 
			  WHERE uc.course_id = c.id AND uc.user_id = $1 AND uc.role IN ('teacher', 'creator')
		  ))
		  AND ($2 = 0 OR c.id = $2)
		  AND c.deleted_at IS NULL
	`

	var total int64
	err = q.QueryRow(ctx, countQuery, teacherID, courseID).Scan(&total)
	if err != nil {
		return nil, 0, nil, fmt.Errorf("%s: count: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	// 3. Select paginated pending homework items
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
		WHERE qaa.is_correct IS NULL
		  AND (c.created_by = $1 OR EXISTS (
			  SELECT 1 FROM users_courses uc 
			  WHERE uc.course_id = c.id AND uc.user_id = $1 AND uc.role IN ('teacher', 'creator')
		  ))
		  AND ($2 = 0 OR c.id = $2)
		  AND c.deleted_at IS NULL
		ORDER BY qa.created_at ASC
		LIMIT $3 OFFSET $4
	`

	rows, err := q.Query(ctx, query, teacherID, courseID, limit, offset)
	if err != nil {
		return nil, 0, nil, fmt.Errorf("%s: query: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}
	defer rows.Close()

	items := make([]domain.PendingHomeworkItem, 0)
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
			return nil, 0, nil, fmt.Errorf("%s: scan: %w: %w", op, errorsAPP.ErrInternalDB, err)
		}

		item.SubmittedAt = submittedAt
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, nil, fmt.Errorf("%s: rows error: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	return items, total, summary, nil
}
