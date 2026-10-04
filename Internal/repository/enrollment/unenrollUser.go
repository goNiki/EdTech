package enrollment

import (
	"context"
	"edtech/internal/infrastructure/txmanager"
	errorsAPP "edtech/pkg/errors"
	"fmt"
)

func (r *repository) UnenrollUser(ctx context.Context, userID, courseID int64) error {
	const op = "repository.enrollment.unenrolluser"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	query := `DELETE FROM users_courses WHERE user_id = $1 AND course_id = $2`

	tag, err := q.Exec(ctx, query, userID, courseID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, errorsAPP.ErrNotEnrolled)
	}

	return nil
}
