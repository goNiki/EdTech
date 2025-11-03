package enrolled

import (
	"context"
	"fmt"
)

func (r *repository) GetRoleUserInCourse(ctx context.Context, userID, courceID int64) (string, error) {

	const op = "repository.enrolled.GetRoleUserInCourse"

	query := `SELECT role FROM users_courses WHERE user_id = $1 AND course_id = $2`

	var role string

	err := r.Pool.QueryRow(ctx, query, userID, courceID).Scan(&role)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return role, nil

}
