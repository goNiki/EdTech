package course

import (
	"context"
	"edtech/internal/domain"
	"fmt"
)

func (r *repository) GetCourseByID(ctx context.Context, id int64) (*domain.Course, error) {

	const op = "repository.course.getcoursebyid"

	query := `SELECT id, title, slug, description, cover_url, created_by, visibility, status, created_at, updated_at FROM courses WHERE id = $1`

	var course domain.Course

	err := r.Pool.QueryRow(ctx, query, id).Scan(&course.Id, &course.Title, &course.Slug, &course.Description, &course.CoverURL, &course.CreatedBy, &course.Visibility, &course.Status, &course.CreatedAt, &course.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &course, nil
}
