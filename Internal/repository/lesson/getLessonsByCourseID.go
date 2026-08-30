package lesson

import (
	"context"
	"edtech/internal/domain"
	"edtech/internal/infrastructure/db"
	"fmt"
)

func (r *repository) GetLessonsByCourseID(ctx context.Context, q db.QueryExecutor, courseID int64) ([]domain.Lesson, error) {
	const op = "repository.lesson.getbycourseid"

	query := `SELECT id, course_id, section_id, title, description, cover_url, content, type, position, duration, is_free, created_at, updated_at, published_at, deleted_at 
		FROM lessons WHERE course_id = $1 AND deleted_at IS NULL ORDER BY position ASC`

	rows, err := q.Query(ctx, query, courseID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var lessons []domain.Lesson
	for rows.Next() {
		var lesson domain.Lesson
		if err := rows.Scan(
			&lesson.ID, &lesson.CourseID, &lesson.SectionID, &lesson.Title, &lesson.Description,
			&lesson.CoverURL, &lesson.Content, &lesson.Type, &lesson.Position, &lesson.Duration,
			&lesson.IsFree, &lesson.CreatedAt, &lesson.UpdatedAt, &lesson.PublishedAt, &lesson.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		lessons = append(lessons, lesson)
	}
	return lessons, nil
}
