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
	if len(entity.QuizSettings) == 0 {
		entity.QuizSettings = []byte(`{"time_limit_minutes":0,"question_time_limit_seconds":0,"max_attempts":0,"passing_score_percent":70,"feedback_mode":"immediate","shuffle_questions":false}`)
	}

	query := `INSERT INTO lessons (course_id, section_id, title, description, cover_url, content, type, position, duration, is_free, status, quiz_settings) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id, created_at, updated_at`

	err := q.QueryRow(ctx, query,
		entity.CourseID, entity.SectionID, entity.Title, entity.Description,
		entity.CoverURL, entity.Content, entity.Type, entity.Position,
		entity.Duration, entity.IsFree, entity.Status, entity.QuizSettings,
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
