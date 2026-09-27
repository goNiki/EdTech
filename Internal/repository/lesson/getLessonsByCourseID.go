package lesson

import (
	"context"
	"fmt"

	"edtech/internal/domain"
	"edtech/internal/infrastructure/db"
	"edtech/internal/repository/models"
	"edtech/internal/repository/models/converter"
	errorsAPP "edtech/pkg/errors"
)

func (r *repository) GetLessonsByCourseID(ctx context.Context, q db.QueryExecutor, courseID int64) ([]domain.Lesson, error) {
	const op = "repository.lesson.getbycourseid"

	query := `SELECT id, course_id, section_id, title, description, cover_url, content, type, position, duration, is_free, status, created_at, updated_at, published_at, deleted_at 
		FROM lessons WHERE course_id = $1 AND deleted_at IS NULL ORDER BY position ASC`

	rows, err := q.Query(ctx, query, courseID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}
	defer rows.Close()

	var lessons []domain.Lesson
	for rows.Next() {
		var m models.Lesson
		if err := rows.Scan(
			&m.ID, &m.CourseID, &m.SectionID, &m.Title, &m.Description,
			&m.CoverURL, &m.Content, &m.Type, &m.Position, &m.Duration,
			&m.IsFree, &m.Status, &m.CreatedAt, &m.UpdatedAt, &m.PublishedAt, &m.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("%s: scan lesson: %w", op, err)
		}
		lessons = append(lessons, *converter.LessonToDomain(&m))
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows error: %w", op, err)
	}

	if lessons == nil {
		lessons = []domain.Lesson{}
	}

	return lessons, nil
}
