package lesson

import (
	"context"
	"fmt"

	"edtech/internal/infrastructure/txmanager"
	errorsAPP "edtech/pkg/errors"
)

func (r *repository) UpdateLessonStatus(ctx context.Context, lessonID int64, status string) error {
	const op = "repository.lesson.UpdateLessonStatus"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	query := `UPDATE lessons SET status = $1, updated_at = NOW() WHERE id = $2 AND deleted_at IS NULL`

	tag, err := q.Exec(ctx, query, status, lessonID)
	if err != nil {
		return fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, errorsAPP.ErrLessonNotFound)
	}
	return nil
}

func (r *repository) UpdateStatusByCourseID(ctx context.Context, courseID int64, status string) error {
	const op = "repository.lesson.UpdateStatusByCourseID"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	query := `UPDATE lessons SET status = $1, updated_at = NOW() WHERE course_id = $2 AND deleted_at IS NULL`

	_, err := q.Exec(ctx, query, status, courseID)
	if err != nil {
		return fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}
	return nil
}

func (r *repository) UpdateStatusBySectionID(ctx context.Context, sectionID int64, status string) error {
	const op = "repository.lesson.UpdateStatusBySectionID"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	query := `UPDATE lessons SET status = $1, updated_at = NOW() WHERE section_id = $2 AND deleted_at IS NULL`

	_, err := q.Exec(ctx, query, status, sectionID)
	if err != nil {
		return fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}
	return nil
}
