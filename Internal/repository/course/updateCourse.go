package course

import (
	"context"
	"edtech/internal/domain"
	errorsAPP "edtech/pkg/errors"
	"fmt"
)

func (r *repository) UpdateCourse(ctx context.Context, course *domain.Course) error {
	const op = "repository.course.updatecourse"

	query := `
		UPDATE courses
		SET title = $1,
			slug = $2,
			description = $3,
			cover_url = $4,
			visibility = $5,
			status = $6,
			updated_at = NOW()
		WHERE id = $7
			AND (
				title IS DISTINCT FROM $1 OR
				slug IS DISTINCT FROM $2 OR
				description IS DISTINCT FROM $3 OR
				cover_url IS DISTINCT FROM $4 OR
				visibility IS DISTINCT FROM $5 OR
				status IS DISTINCT FROM $6
			)`

	cmtTag, err := r.Pool.Exec(ctx, query, course.Title, course.Slug, course.Description, course.CoverURL, course.Visibility, course.Status, course.Id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if cmtTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, errorsAPP.ErrNothingToUpdate)
	}

	return nil
}
