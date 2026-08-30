package course

import (
	"context"
	"edtech/internal/infrastructure/db"
	errorsAPP "edtech/pkg/errors"
	"fmt"
)

func (r *repository) IncrementEnrolledCount(ctx context.Context, q db.QueryExecutor, courseID int64) error {
	const op = "repository.course.IncrementEnrolledCount"

	query := `UPDATE courses SET enrolled_count = enrolled_count + 1 WHERE id = $1`

	tag, err := q.Exec(ctx, query, courseID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, errorsAPP.ErrNotFoundCourse)
	}

	return nil
}

func (r *repository) DecrementEnrolledCount(ctx context.Context, q db.QueryExecutor, courseID int64) error {
	const op = "repository.course.DecrementEnrolledCount"

	query := `UPDATE courses SET enrolled_count = enrolled_count - 1 WHERE id = $1`

	tag, err := q.Exec(ctx, query, courseID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, errorsAPP.ErrNotFoundCourse)
	}

	return nil
}
