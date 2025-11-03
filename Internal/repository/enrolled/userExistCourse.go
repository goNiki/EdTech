package enrolled

import (
	"context"
	"fmt"
)

func (r *repository) UserExistCourse(ctx context.Context, userID, courseid int64) (bool, error) {

	const op = "repository.enrolled.userexistcourse"

	query := `SELECT 1 FROM users_courses WHERE user_id = $1 AND course_id = $2 LIMIT 1`

	var exist int

	err := r.Pool.QueryRow(ctx, query, userID, courseid).Scan(&exist)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return true, nil
}
