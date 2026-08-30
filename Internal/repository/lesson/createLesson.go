package lesson

import (
	"context"
	"edtech/internal/domain"
	"edtech/internal/infrastructure/db"
	"fmt"
)

func (r *repository) CreateLesson(ctx context.Context, q db.QueryExecutor, lesson *domain.Lesson) error {
	const op = "repository.lesson.create"

	query := `INSERT INTO lessons (course_id, section_id, title, description, cover_url, content, type, position, duration, is_free) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`

	err := q.QueryRow(ctx, query, lesson.CourseID, lesson.SectionID, lesson.Title, lesson.Description,
		lesson.CoverURL, lesson.Content, lesson.Type, lesson.Position, lesson.Duration, lesson.IsFree,
	).Scan(&lesson.ID)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
