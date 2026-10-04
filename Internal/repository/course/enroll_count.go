package course

import (
	"context"
	"edtech/internal/infrastructure/txmanager"
	errorsAPP "edtech/pkg/errors"
	"fmt"
)

func (r *repository) IncrementEnrolledCount(ctx context.Context, courseID int64) error {
	const op = "repository.course.IncrementEnrolledCount"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

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

func (r *repository) DecrementEnrolledCount(ctx context.Context, courseID int64) error {
	const op = "repository.course.DecrementEnrolledCount"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	query := `UPDATE courses SET enrolled_count = GREATEST(enrolled_count - 1, 0) WHERE id = $1`

	tag, err := q.Exec(ctx, query, courseID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, errorsAPP.ErrNotFoundCourse)
	}

	return nil
}
