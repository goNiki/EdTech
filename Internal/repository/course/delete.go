package course

import (
	"context"
	"fmt"

	"edtech/internal/infrastructure/db"
	errorsAPP "edtech/pkg/errors"
)

func (r *repository) DeleteCourse(ctx context.Context, q db.QueryExecutor, courseID int64) error {
	const op = "repository.course.DeleteCourse"

	query := `
		UPDATE courses 
		SET 
			deleted_at = NOW(),
			status = 'archived',
			slug = 'deleted_' || id::text || '_' || slug,
			updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	tag, err := q.Exec(ctx, query, courseID)
	if err != nil {
		return fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, errorsAPP.ErrCourseNotFound)
	}

	return nil
}
