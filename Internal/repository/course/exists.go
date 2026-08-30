package course

import (
	"context"
	"edtech/internal/infrastructure/db"
	errorsAPP "edtech/pkg/errors"
	"fmt"
)

func (r *repository) ExistingBySlug(ctx context.Context, q db.QueryExecutor, slug string) (bool, error) {
	const op = "repository.course.ExistingBySlug"

	query := `
		SELECT EXISTS(
			SELECT 1
			FROM courses
			WHERE slug = $1 AND deleted_at IS NULL
		)
	`

	var exists bool

	if err := q.QueryRow(ctx, query, slug).Scan(&exists); err != nil {
		return false, fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	return exists, nil

}
