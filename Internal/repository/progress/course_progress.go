package progress

import (
	"context"
	"errors"
	"fmt"

	"edtech/internal/domain"
	"edtech/internal/infrastructure/db"
	"edtech/internal/repository/models"
	repoconverter "edtech/internal/repository/models/converter"
	errorsAPP "edtech/pkg/errors"

	"github.com/jackc/pgx/v5"
)

func (r *repository) CreateCourseProgress(ctx context.Context, q db.QueryExecutor, progress *domain.CourseProgress) error {
	const op = "repository.progress.CreateCourseProgress"

	query := `
		INSERT INTO course_progress (user_id, course_id, completed_lessons, total_lessons, progress_percentage, total_watch_time, average_score, started_at, last_accessed_at, completed_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, COALESCE($8, NOW()), NOW(), $9)
		ON CONFLICT (user_id, course_id) DO UPDATE
		SET completed_lessons = EXCLUDED.completed_lessons,
		    total_lessons = EXCLUDED.total_lessons,
		    progress_percentage = EXCLUDED.progress_percentage,
		    last_accessed_at = NOW(),
		    completed_at = CASE WHEN EXCLUDED.progress_percentage >= 100 AND course_progress.completed_at IS NULL THEN NOW() ELSE course_progress.completed_at END
		RETURNING id, last_accessed_at
	`

	err := q.QueryRow(
		ctx,
		query,
		progress.UserID,
		progress.CourseID,
		progress.CompletedLess,
		progress.TotalLessons,
		progress.Percent,
		progress.TotalWatchTime,
		progress.AverageScore,
		progress.StartedAt,
		progress.CompletedAt,
	).Scan(&progress.ID, &progress.LastAccessedAt)

	if err != nil {
		return fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	return nil
}

func (r *repository) UpsertCourseProgress(ctx context.Context, q db.QueryExecutor, userID, courseID int64, completedLessons, totalLessons int, percentage float64) error {
	const op = "repository.progress.UpsertCourseProgress"

	query := `
		INSERT INTO course_progress (user_id, course_id, completed_lessons, total_lessons, progress_percentage, started_at, last_accessed_at, completed_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW(), CASE WHEN $5 >= 100 THEN NOW() ELSE NULL END)
		ON CONFLICT (user_id, course_id) DO UPDATE
		SET completed_lessons = EXCLUDED.completed_lessons,
		    total_lessons = EXCLUDED.total_lessons,
		    progress_percentage = EXCLUDED.progress_percentage,
		    last_accessed_at = NOW(),
		    completed_at = CASE WHEN EXCLUDED.progress_percentage >= 100 AND course_progress.completed_at IS NULL THEN NOW() ELSE course_progress.completed_at END
	`

	_, err := q.Exec(ctx, query, userID, courseID, completedLessons, totalLessons, percentage)
	if err != nil {
		return fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	return nil
}

func (r *repository) GetCourseProgress(ctx context.Context, q db.QueryExecutor, userID, courseID int64) (*domain.CourseProgress, error) {
	const op = "repository.progress.GetCourseProgress"

	query := `
		SELECT id, user_id, course_id, completed_lessons, total_lessons, progress_percentage, total_watch_time, average_score, started_at, last_accessed_at, completed_at
		FROM course_progress
		WHERE user_id = $1 AND course_id = $2
	`

	var cp models.CourseProgress
	err := q.QueryRow(ctx, query, userID, courseID).Scan(
		&cp.ID,
		&cp.UserID,
		&cp.CourseID,
		&cp.CompletedLess,
		&cp.TotalLessons,
		&cp.Percent,
		&cp.TotalWatchTime,
		&cp.AverageScore,
		&cp.StartedAt,
		&cp.LastAccessedAt,
		&cp.CompletedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, errorsAPP.ErrCourseProgressNotFound)
		}
		return nil, fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	return repoconverter.CourseProgressToDomain(&cp), nil
}

func (r *repository) GetAllLessonProgressByCourse(ctx context.Context, q db.QueryExecutor, userID, courseID int64) ([]domain.LessonProgress, error) {
	const op = "repository.progress.GetAllLessonProgressByCourse"

	query := `
		SELECT id, user_id, lesson_id, course_id, status, score, watch_time, last_position, started_at, completed_at, last_accessed_at
		FROM lesson_progress
		WHERE user_id = $1 AND course_id = $2
		ORDER BY id ASC
	`

	rows, err := q.Query(ctx, query, userID, courseID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}
	defer rows.Close()

	var list []domain.LessonProgress
	for rows.Next() {
		var progress models.LessonProgress
		var rawScore *float64

		if err := rows.Scan(
			&progress.ID,
			&progress.UserID,
			&progress.LessonID,
			&progress.CourseID,
			&progress.Status,
			&rawScore,
			&progress.TimeSpent,
			&progress.LastPos,
			&progress.StartedAt,
			&progress.CompletedAt,
			&progress.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		if rawScore != nil {
			sc := int(*rawScore)
			progress.Score = &sc
		}
		list = append(list, *repoconverter.LessonProgressToDomain(&progress))
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if list == nil {
		list = []domain.LessonProgress{}
	}

	return list, nil
}
