package course

import (
	"context"
	"edtech/internal/domain"
	errorsAPP "edtech/pkg/errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (r *repository) GetCourseBySlug(ctx context.Context, slug string) (*domain.Course, error) {
	const op = "repository.course.getcoursebyslug"

	query := `SELECT id, title, slug, description, cover_url, created_by, visibility, status, created_at, updated_at FROM courses WHERE slug = $1`

	var course domain.Course

	err := r.Pool.QueryRow(ctx, query, slug).Scan(&course.Id, &course.Title, &course.Slug, &course.Description, &course.CoverURL, &course.CreatedBy, &course.Visibility, &course.Status, &course.CreatedAt, &course.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("%s: %w", op, errorsAPP.ErrNotFoundCourse)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &course, nil
}
