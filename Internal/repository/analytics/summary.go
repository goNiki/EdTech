package analytics

import (
	"context"
	"fmt"

	"edtech/internal/domain"
	"edtech/internal/infrastructure/txmanager"
	errorsAPP "edtech/pkg/errors"
)

func (r *repositoryAnalytics) GetCourseAnalyticsSummary(ctx context.Context, courseID int64) (domain.CourseAnalyticsSummary, error) {
	const op = "repository.analytics.GetCourseAnalyticsSummary"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	query := `
		WITH student_count AS (
			SELECT COUNT(*) AS total_students 
			FROM users_courses 
			WHERE course_id = $1 AND role = 'student'
		),
		course_total_lessons AS (
			SELECT COUNT(*) AS total_lessons 
			FROM lessons 
			WHERE course_id = $1 AND deleted_at IS NULL
		),
		progress_stats AS (
			SELECT 
				COALESCE(AVG(
					CASE 
						WHEN ctl.total_lessons > 0 THEN (sp.completed_count::float / ctl.total_lessons::float) * 100.0 
						ELSE 0.0 
					END
				), 0.0) AS avg_progress_percent,
				COALESCE(AVG(sp.avg_score), 0.0) AS avg_score
			FROM (
				SELECT 
					lp.user_id,
					COUNT(*) FILTER (WHERE lp.status = 'completed') AS completed_count,
					AVG(lp.score) FILTER (WHERE lp.score > 0) AS avg_score
				FROM lesson_progress lp
				WHERE lp.course_id = $1
				GROUP BY lp.user_id
			) sp
			CROSS JOIN course_total_lessons ctl
		),
		pending_hw AS (
			SELECT COUNT(*) AS pending_count
			FROM quiz_attempt_answers qaa
			JOIN quiz_attempts qa ON qaa.attempt_id = qa.id
			JOIN quizzes qz ON qa.quiz_id = qz.id
			WHERE qz.course_id = $1 AND qaa.is_correct IS NULL
		)
		SELECT 
			sc.total_students,
			COALESCE(ps.avg_progress_percent, 0.0),
			COALESCE(ps.avg_score, 0.0),
			ph.pending_count
		FROM student_count sc
		CROSS JOIN progress_stats ps
		CROSS JOIN pending_hw ph
	`

	var summary domain.CourseAnalyticsSummary
	err := q.QueryRow(ctx, query, courseID).Scan(
		&summary.TotalStudents,
		&summary.AvgProgressPercent,
		&summary.AvgScore,
		&summary.PendingHomeworksCount,
	)
	if err != nil {
		return domain.CourseAnalyticsSummary{}, fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	return summary, nil
}
