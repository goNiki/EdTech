package lesson

import (
	"context"
	"edtech/internal/domain"
	errorsAPP "edtech/pkg/errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (r *repository) GetLessonByID(ctx context.Context, id int) (*domain.Lesson, error) {
	const op = "repositiry.couse.lessonrepo.GetLessonByid"

	query := `SELECT id, course_id, title, description, cover_url, content, position, create_at, update_at FROM lessons WHERE id = $1`

	var lesson domain.Lesson
	err := r.Pool.QueryRow(ctx, query, id).Scan(&lesson.ID, &lesson.CourseID, &lesson.Title, &lesson.Description, &lesson.CoverURL, &lesson.Content, &lesson.Position, &lesson.CreatedAt, &lesson.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("%s: %w", op, errorsAPP.ErrNotFoundLesson)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &lesson, nil

}
