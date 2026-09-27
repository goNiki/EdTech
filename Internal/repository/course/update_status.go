package course

import (
	"context"
	"fmt"

	"edtech/internal/infrastructure/db"
	errorsAPP "edtech/pkg/errors"
)

func (r *repository) UpdateCourseStatus(ctx context.Context, q db.QueryExecutor, courseID int64, status string) error {
	const op = "repository.course.UpdateCourseStatus"

	query := `UPDATE courses SET status = $1, updated_at = NOW() WHERE id = $2 AND deleted_at IS NULL`

	tag, err := q.Exec(ctx, query, status, courseID)
	if err != nil {
		return fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, errorsAPP.ErrCourseNotFound)
	}

	return nil
}
