package enrollment

import (
	"context"
	"edtech/internal/infrastructure/db"
	errorsAPP "edtech/pkg/errors"
	"fmt"
)

func (r *repository) UnenrollUser(ctx context.Context, q db.QueryExecutor, userID, courseID int64) error {
	const op = "repository.enrollment.unenrolluser"

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
