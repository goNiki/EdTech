package course

import (
	"context"
	"edtech/internal/domain"
	errorsAPP "edtech/pkg/errors"
	"fmt"
)

func (r *repository) CreateCourse(ctx context.Context, course *domain.Course) (int64, error) {
	const op = "repository.course.createcourse"

	query := `INSERT INTO courses (title, slug, description, cover_url, created_by, visibility, status) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id `

	err := r.Pool.QueryRow(ctx, query, course.Title, course.Slug, course.Description, course.CoverURL, course.CreatedBy, course.Visibility, course.Status).Scan(&course.Id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	query = `INSERT INTO users_courses (user_id, course_id, role) VALUES ($1, $2, $3)`

	role := "teacher"

	cmgTag, err := r.Pool.Exec(ctx, query, course.CreatedBy, course.Id, role)

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	if cmgTag.RowsAffected() == 0 {
		return 0, fmt.Errorf("%s: %w", op, errorsAPP.ErrEnrolledByCreated)
	}

	return int64(course.Id), nil
}
