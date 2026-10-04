package course

import (
	"context"
	"fmt"

	"edtech/internal/infrastructure/txmanager"
	errorsAPP "edtech/pkg/errors"
)

func (r *repository) UpdateCourseRatingStats(ctx context.Context, courseID int64, rating float64, reviewsCount int) error {
	const op = "repository.course.UpdateCourseRatingStats"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	query := `UPDATE courses SET rating = $1, reviews_count = $2 WHERE id = $3 AND deleted_at IS NULL`

	tag, err := q.Exec(ctx, query, rating, reviewsCount, courseID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, errorsAPP.ErrCourseNotFound)
	}

	return nil
}
