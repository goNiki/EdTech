package analytics

import (
	"context"
	"fmt"
	"math"
	"time"

	"edtech/internal/domain"
	"edtech/internal/infrastructure/txmanager"
	errorsAPP "edtech/pkg/errors"
)

func (r *repositoryAnalytics) GetCourseGradebook(ctx context.Context, courseID int64) ([]domain.GradebookRecord, error) {
	const op = "repository.analytics.GetCourseGradebook"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

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
		)
		SELECT 
			u.id,
			COALESCE(NULLIF(TRIM(u.first_name || ' ' || u.last_name), ''), u.username) AS student_name,
			u.email,
			uc.enrolled_at,
			COALESCE(sp.completed_lessons, 0) AS completed_lessons,
			COALESCE(cs.total_lessons, 0) AS total_lessons,
			COALESCE(sp.avg_score, 0) AS avg_score,
			cert.certificate_code
		FROM users_courses uc
		JOIN users u ON uc.user_id = u.id
		CROSS JOIN course_stats cs
		LEFT JOIN student_progress sp ON sp.user_id = u.id
		LEFT JOIN certificates cert ON cert.user_id = u.id AND cert.course_id = $1
		WHERE uc.course_id = $1 AND uc.role = 'student' AND u.deleted_at IS NULL
		ORDER BY uc.enrolled_at ASC, u.id ASC
	`

	rows, err := q.Query(ctx, query, courseID)
	if err != nil {
		return nil, fmt.Errorf("%s: query: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}
	defer rows.Close()

	var records []domain.GradebookRecord
	for rows.Next() {
		var (
			rec      domain.GradebookRecord
			enrolled time.Time
			avgScore float64
			certCode *string
		)

		if err := rows.Scan(
			&rec.UserID,
			&rec.StudentName,
			&rec.Email,
			&enrolled,
			&rec.CompletedLessons,
			&rec.TotalLessons,
			&avgScore,
			&certCode,
		); err != nil {
			return nil, fmt.Errorf("%s: scan: %w: %w", op, errorsAPP.ErrInternalDB, err)
		}

		rec.EnrolledAt = enrolled
		rec.CertificateCode = certCode

		if rec.TotalLessons > 0 {
			pct := float64(rec.CompletedLessons*100) / float64(rec.TotalLessons)
			if pct > 100.0 {
				pct = 100.0
			}
			rec.ProgressPercent = int(math.Round(pct))
		} else {
			rec.ProgressPercent = 0
		}

		rec.AverageScore = int(math.Round(avgScore))

		if rec.ProgressPercent >= 100 && rec.TotalLessons > 0 {
			rec.Status = "Завершен"
		} else {
			rec.Status = "В процессе"
		}

		records = append(records, rec)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	if records == nil {
		records = []domain.GradebookRecord{}
	}

	return records, nil
}
