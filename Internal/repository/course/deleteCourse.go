package course

import (
	"context"
	"fmt"
)

func (r *repository) DeleteCourse(ctx context.Context, courseID int64) error {
	const op = "repository.course.deletecourse"

	query := `DELETE FROM courses WHERE id = $1`

	_, err := r.Pool.Exec(ctx, query, courseID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil

}
