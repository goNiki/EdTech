package course

import (
	"context"
	"fmt"
)

func (r *repository) CountCourse(ctx context.Context) (int, error) {

	const op = "repository.course.countcourse"

	query := `SELECT count(*) FROM courses WHERE status = 'published'`

	var total int

	err := r.Pool.QueryRow(ctx, query).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return total, nil

}
