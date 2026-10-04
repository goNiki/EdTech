package enrollment

import (
	"context"
	"fmt"

	"edtech/internal/infrastructure/txmanager"
)

func (r *repository) GetRoleUserInCourse(ctx context.Context, userID, courceID int64) (string, error) {
	const op = "repository.enrolled.GetRoleUserInCourse"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	query := `SELECT role FROM users_courses WHERE user_id = $1 AND course_id = $2`

	var role string

	err := q.QueryRow(ctx, query, userID, courceID).Scan(&role)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return "", nil
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return role, nil
}
