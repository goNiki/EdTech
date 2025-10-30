package lesson

import (
	"context"
	"edtech/internal/domain"
	"fmt"
)

func (r *repository) CreateLesson(ctx context.Context, lesson *domain.Lesson) error {
	const op = "repository.course.lesson.createlesson"

	query := `
		INSERT INTO lessons (course_id, title, description, cover_url, content, position) VALUES ($1, $2, $3, $4, %5, %6) RETURNING id`

	err := r.Pool.QueryRow(ctx, query, lesson.CourseID, lesson.Title, lesson.Description, lesson.CoverURL, lesson.Content, lesson.Position).Scan(&lesson.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil

}
