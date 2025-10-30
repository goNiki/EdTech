package course

import (
	"context"
	"edtech/internal/domain"
	"fmt"
)

func (r *repository) PublishCourse(ctx context.Context, course *domain.Course) error {
	const op = "repository.course.publishcourse"

	query := `UPDATE courses SET status = $1 , updated_at = $2 WHERE id = $3`

	_, err := r.Pool.Exec(ctx, query, course.Status, course.UpdatedAt, course.Id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil

}
