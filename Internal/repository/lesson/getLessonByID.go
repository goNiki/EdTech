package lesson

import (
	"context"
	"errors"
	"fmt"

	"edtech/internal/domain"
	"edtech/internal/infrastructure/db"
	"edtech/internal/repository/models"
	"edtech/internal/repository/models/converter"
	errorsAPP "edtech/pkg/errors"

	"github.com/jackc/pgx/v5"
)

func (r *repository) GetLessonByID(ctx context.Context, q db.QueryExecutor, id int64) (*domain.Lesson, error) {
	const op = "repository.lesson.getbyid"

	query := `SELECT id, course_id, section_id, title, description, cover_url, content, type, position, duration, is_free, status, created_at, updated_at, published_at, deleted_at 
		FROM lessons WHERE id = $1 AND deleted_at IS NULL`

	var m models.Lesson
	err := q.QueryRow(ctx, query, id).Scan(
		&m.ID, &m.CourseID, &m.SectionID, &m.Title, &m.Description,
		&m.CoverURL, &m.Content, &m.Type, &m.Position, &m.Duration,
		&m.IsFree, &m.Status, &m.CreatedAt, &m.UpdatedAt, &m.PublishedAt, &m.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, errorsAPP.ErrLessonNotFound)
		}
		return nil, fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	return converter.LessonToDomain(&m), nil
}
