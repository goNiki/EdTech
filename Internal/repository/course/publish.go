package course

import (
	"context"
	"fmt"

	"edtech/internal/domain"
	"edtech/internal/infrastructure/txmanager"
	errorsAPP "edtech/pkg/errors"
)

func (r *repository) PublishCourse(ctx context.Context, course *domain.Course) error {
	const op = "repository.course.publishcourse"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	query := `
		UPDATE courses 
		SET 
			status = $1, 
			updated_at = $2, 
			published_at = $3 
		WHERE id = $4 AND deleted_at IS NULL
	`

	tag, err := q.Exec(
		ctx,
		query,
		course.Status,
		course.UpdatedAt,
		course.PublishedAt,
		course.Id,
	)
	if err != nil {
		return fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, errorsAPP.ErrNothingToUpdate)
	}

	return nil
}
