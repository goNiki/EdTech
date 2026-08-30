package enrollment

import (
	"context"
	"edtech/internal/infrastructure/db"
	errorsAPP "edtech/pkg/errors"
	"fmt"
)

func (r *repository) ChangeUserRole(ctx context.Context, q db.QueryExecutor, courseID int64, targetUserID int64, newRole string) error {
	const op = "repository.enrollment.changeuserrole"

	query := `UPDATE users_courses SET role = $1 WHERE user_id = $2 AND course_id = $3`

	tag, err := q.Exec(ctx, query, newRole, targetUserID, courseID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, errorsAPP.ErrNotEnrolled)
	}

	return nil
}
