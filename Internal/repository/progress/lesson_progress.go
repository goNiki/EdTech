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

func (r *repository) CreateLessonProgress(ctx context.Context, q db.QueryExecutor, progress *domain.LessonProgress) error {
	const op = "repository.progress.CreateLessonProgress"

	status := progress.Status
	if status == "" {
		status = domain.ProgressStatusInProgress
	}

	query := `
		INSERT INTO lesson_progress (user_id, lesson_id, course_id, status, last_position, watch_time, score, started_at, completed_at, last_accessed_at)
		VALUES ($1, $2, $3, $4::progress_status, $5, $6, $7, $8, $9, NOW())
		ON CONFLICT (user_id, lesson_id) DO UPDATE
		SET status = EXCLUDED.status,
		    completed_at = COALESCE(EXCLUDED.completed_at, lesson_progress.completed_at),
		    last_accessed_at = NOW()
		RETURNING id, last_accessed_at
	`

	err := q.QueryRow(
		ctx,
		query,
		progress.UserID,
		progress.LessonID,
		progress.CourseID,
		string(status),
		progress.LastPos,
		progress.TimeSpent,
		progress.Score,
		progress.StartedAt,
		progress.CompletedAt,
	).Scan(&progress.ID, &progress.UpdatedAt)

	if err != nil {
		return fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	return nil
}

func (r *repository) GetLessonProgress(ctx context.Context, q db.QueryExecutor, userID, lessonID int64) (*domain.LessonProgress, error) {
	const op = "repository.progress.GetLessonProgress"

	query := `
		SELECT id, user_id, lesson_id, course_id, status, score, watch_time, last_position, started_at, completed_at, last_accessed_at
		FROM lesson_progress
		WHERE user_id = $1 AND lesson_id = $2
	`

	var progress models.LessonProgress
	var rawScore *float64

	err := q.QueryRow(ctx, query, userID, lessonID).Scan(
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
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, errorsAPP.ErrLessonProgressNotFound)
		}
		return nil, fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	if rawScore != nil {
		sc := int(*rawScore)
		progress.Score = &sc
	}

	return repoconverter.LessonProgressToDomain(&progress), nil
}

func (r *repository) UpdateLessonProgressTime(ctx context.Context, q db.QueryExecutor, userID, lessonID int64, additionalTime int, lastPos int) error {
	const op = "repository.progress.UpdateLessonProgressTime"

	query := `
		UPDATE lesson_progress 
		SET watch_time = watch_time + $1, 
		    last_position = $2, 
		    last_accessed_at = NOW() 
		WHERE user_id = $3 AND lesson_id = $4
	`

	tag, err := q.Exec(ctx, query, additionalTime, lastPos, userID, lessonID)
	if err != nil {
		return fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, errorsAPP.ErrLessonProgressNotFound)
	}

	return nil
}

func (r *repository) UpdateLessonProgressStatus(ctx context.Context, q db.QueryExecutor, userID, lessonID int64, status domain.ProgressStatus) error {
	const op = "repository.progress.UpdateLessonProgressStatus"

	query := `
		UPDATE lesson_progress 
		SET status = $1::progress_status,
		    started_at = CASE WHEN $1::text = 'in_progress' AND started_at IS NULL THEN NOW() ELSE started_at END,
		    completed_at = CASE WHEN $1::text = 'completed' AND completed_at IS NULL THEN NOW() ELSE completed_at END,
		    last_accessed_at = NOW()
		WHERE user_id = $2 AND lesson_id = $3
	`

	tag, err := q.Exec(ctx, query, string(status), userID, lessonID)
	if err != nil {
		return fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, errorsAPP.ErrLessonProgressNotFound)
	}

	return nil
}
