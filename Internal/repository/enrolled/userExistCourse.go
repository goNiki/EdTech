package enrolled

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (r *repository) UserExistCourse(ctx context.Context, userID, courseid int) (bool, error) {

	const op = "repository.enrolled.userexistcourse"

	query := `SELECT 1 FROM users_courses WHERE user_id = $1 AND course_id = $2 LIMIT 1`

	var exist int

	err := r.Pool.QueryRow(ctx, query, userID, courseid).Scan(&exist)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return true, nil
}
