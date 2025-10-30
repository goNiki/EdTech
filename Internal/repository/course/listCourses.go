package course

import (
	"context"
	"edtech/internal/domain"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (r *repository) ListCourses(ctx context.Context, pageSize int, offset int) ([]domain.Course, error) {

	const op = "repository.course.listcourses"

	query := `SELECT id, title, slug, description, cover_url, created_by, visibility, status,  created_at, updated_at FROM courses WHERE status = 'published' ORDER BY create_at DESC LIMIT $1 OFFSET $2`

	rows, err := r.Pool.Query(ctx, query, pageSize, offset)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer rows.Close()
	var courses []domain.Course

	for rows.Next() {
		var c domain.Course

		if err := rows.Scan(&c.Id, &c.Title, &c.Slug, &c.Description, &c.CoverURL, &c.CreatedBy, &c.Visibility, &c.Status, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		courses = append(courses, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if len(courses) == 0 {
		return nil, fmt.Errorf("%s: %w", op, pgx.ErrNoRows)
	}

	return courses, nil

}
