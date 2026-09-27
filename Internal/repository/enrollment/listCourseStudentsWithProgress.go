package enrollment

import (
	"context"
	"fmt"
	"time"

	"edtech/internal/domain"
	"edtech/internal/infrastructure/db"
	errorsAPP "edtech/pkg/errors"
)

func (r *repository) ListCourseStudentsWithProgress(ctx context.Context, q db.QueryExecutor, courseID int64, limit, offset int64) ([]domain.CourseStudentItem, int64, error) {
	const op = "repository.enrollment.ListCourseStudentsWithProgress"

	countQuery := `
		SELECT COUNT(*) 
		FROM users_courses uc
		JOIN users u ON uc.user_id = u.id
		WHERE uc.course_id = $1 AND uc.role = 'student' AND u.deleted_at IS NULL
	`
	var total int64
	err := q.QueryRow(ctx, countQuery, courseID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: count: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	query := `
		WITH course_stats AS (
			SELECT COUNT(*) AS total_lessons
			FROM lessons 
			WHERE course_id = $1 AND deleted_at IS NULL
		),
		student_progress AS (
			SELECT 
				lp.user_id,
				COUNT(*) FILTER (WHERE lp.status = 'completed') AS completed_lessons,
				COALESCE(AVG(lp.score) FILTER (WHERE lp.score > 0), 0) AS avg_score
			FROM lesson_progress lp
			WHERE lp.course_id = $1
			GROUP BY lp.user_id
		),
		pending_stats AS (
			SELECT 
				qa.user_id,
				COUNT(*) AS pending_count
			FROM quiz_attempt_answers qaa
			JOIN quiz_attempts qa ON qaa.attempt_id = qa.id
			JOIN quizzes qz ON qa.quiz_id = qz.id
			WHERE qz.course_id = $1 AND qaa.is_correct IS NULL
			GROUP BY qa.user_id
		)
		SELECT 
			u.id,
			COALESCE(u.first_name, '') AS first_name,
			COALESCE(u.last_name, '') AS last_name,
			u.username,
			u.email,
			u.avatar_url,
			uc.role,
			uc.enrolled_at,
			COALESCE(sp.completed_lessons, 0) AS completed_lessons,
			COALESCE(cs.total_lessons, 0) AS total_lessons,
			COALESCE(sp.avg_score, 0) AS average_score,
			(COALESCE(pend.pending_count, 0) > 0) AS has_pending_homeworks
		FROM users_courses uc
		JOIN users u ON uc.user_id = u.id
		CROSS JOIN course_stats cs
		LEFT JOIN student_progress sp ON sp.user_id = u.id
		LEFT JOIN pending_stats pend ON pend.user_id = u.id
		WHERE uc.course_id = $1 AND uc.role = 'student' AND u.deleted_at IS NULL
		ORDER BY uc.enrolled_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := q.Query(ctx, query, courseID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: query: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}
	defer rows.Close()

	var students []domain.CourseStudentItem
	for rows.Next() {
		var item domain.CourseStudentItem
		var enrolledAt time.Time

		if err := rows.Scan(
			&item.UserID,
			&item.FirstName,
			&item.LastName,
			&item.Username,
			&item.Email,
			&item.AvatarURL,
			&item.Role,
			&enrolledAt,
			&item.CompletedLessons,
			&item.TotalLessons,
			&item.AverageScore,
			&item.HasPendingHomeworks,
		); err != nil {
			return nil, 0, fmt.Errorf("%s: scan: %w: %w", op, errorsAPP.ErrInternalDB, err)
		}

		item.EnrolledAt = enrolledAt
		if item.TotalLessons > 0 {
			item.ProgressPercentage = float64(item.CompletedLessons*100) / float64(item.TotalLessons)
			if item.ProgressPercentage > 100.0 {
				item.ProgressPercentage = 100.0
			}
		} else {
			item.ProgressPercentage = 0.0
		}

		students = append(students, item)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("%s: rows error: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	if students == nil {
		students = []domain.CourseStudentItem{}
	}

	return students, total, nil
}
