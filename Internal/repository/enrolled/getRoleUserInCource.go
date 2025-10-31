package enrolled

import "context"

func (r *repository) GetRoleUserInCource(ctx context.Context, userID, courceID int64) (string, error) {

	query := `SELECT role FROM users_courses WHERE user_id = $1 AND course_id = $2`

	var role string

	err := r.Pool.QueryRow(ctx, query, userID, courceID).Scan(&role)
	if err != nil {
		return "", err
	}

	return role, nil

}
