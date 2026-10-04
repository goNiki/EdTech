package lesson

import (
	"context"
	"fmt"

	"edtech/internal/domain"
	"edtech/internal/infrastructure/txmanager"
	"edtech/internal/repository/models/converter"
	errorsAPP "edtech/pkg/errors"
)

func (r *repository) CreateLesson(ctx context.Context, lesson *domain.Lesson) error {
	const op = "repository.lesson.create"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	entity := converter.LessonToEntity(lesson)
	if entity.Status == "" {
		entity.Status = domain.StatusDraft
	}

	query := `INSERT INTO lessons (course_id, section_id, title, description, cover_url, content, type, position, duration, is_free, status) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id, created_at, updated_at`

	err := q.QueryRow(ctx, query,
		entity.CourseID, entity.SectionID, entity.Title, entity.Description,
		entity.CoverURL, entity.Content, entity.Type, entity.Position,
		entity.Duration, entity.IsFree, entity.Status,
	).Scan(&entity.ID, &entity.CreatedAt, &entity.UpdatedAt)

	if err != nil {
		return fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	lesson.ID = entity.ID
	lesson.CreatedAt = entity.CreatedAt
	lesson.UpdatedAt = entity.UpdatedAt
	lesson.Status = entity.Status
	return nil
}
