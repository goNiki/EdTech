package lesson

import (
	"context"
	"edtech/internal/domain"
	"edtech/internal/infrastructure/db"
	"fmt"
)

func (r *repository) GetLessonByID(ctx context.Context, q db.QueryExecutor, id int64) (*domain.Lesson, error) {
	const op = "repository.lesson.getbyid"

	query := `SELECT id, course_id, section_id, title, description, cover_url, content, type, position, duration, is_free, created_at, updated_at, published_at, deleted_at 
		FROM lessons WHERE id = $1 AND deleted_at IS NULL`

	var lesson domain.Lesson
	err := q.QueryRow(ctx, query, id).Scan(
		&lesson.ID, &lesson.CourseID, &lesson.SectionID, &lesson.Title, &lesson.Description,
		&lesson.CoverURL, &lesson.Content, &lesson.Type, &lesson.Position, &lesson.Duration,
		&lesson.IsFree, &lesson.CreatedAt, &lesson.UpdatedAt, &lesson.PublishedAt, &lesson.DeletedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &lesson, nil
}
