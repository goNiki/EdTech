package course

import (
	"context"
	"fmt"

	"edtech/internal/domain"
	"edtech/internal/infrastructure/db"
	errorsAPP "edtech/pkg/errors"
)

func (r *repository) ArchiveCourse(ctx context.Context, q db.QueryExecutor, course *domain.Course) error {
	const op = "repository.course.ArchiveCourse"

	query := `
		UPDATE courses 
		SET 
			status = $1, 
			updated_at = $2, 
			archived_at = $3 
		WHERE id = $4 AND deleted_at IS NULL
	`

	tag, err := q.Exec(
		ctx,
		query,
		course.Status,
		course.UpdatedAt,
		course.ArchivedAt,
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
