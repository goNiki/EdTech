package enrollment

import (
	"context"
	"fmt"

	"edtech/internal/infrastructure/db"
)

func (r *repository) UserExistCourse(ctx context.Context, q db.QueryExecutor, userID, courseid int64) (bool, error) {
	const op = "repository.enrolled.userexistcourse"

	query := `
		SELECT EXISTS(
			SELECT 1 
			FROM users_courses 
			WHERE user_id = $1 AND course_id = $2
		)
	`

	var exists bool

	err := q.QueryRow(ctx, query, userID, courseid).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return exists, nil
}
