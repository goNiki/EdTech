package enrollment

import (
	"context"
	"fmt"

	"edtech/internal/infrastructure/txmanager"
)

func (r *repository) UserExistCourse(ctx context.Context, userID, courseid int64) (bool, error) {
	const op = "repository.enrolled.userexistcourse"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

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
