package lesson

import (
	"context"
	"fmt"

	"edtech/internal/domain"
	"edtech/internal/infrastructure/txmanager"
	"edtech/internal/repository/models/converter"
	errorsAPP "edtech/pkg/errors"
)

func (r *repository) UpdateLesson(ctx context.Context, lesson *domain.Lesson) error {
	const op = "repository.lesson.update"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	entity := converter.LessonToEntity(lesson)
	status := entity.Status
	if status == "" {
		status = domain.StatusDraft
	}
	if len(entity.QuizSettings) == 0 {
		entity.QuizSettings = []byte(`{"time_limit_minutes":0,"question_time_limit_seconds":0,"max_attempts":0,"passing_score_percent":70,"feedback_mode":"immediate","shuffle_questions":false}`)
	}

	query := `UPDATE lessons SET course_id = $1, section_id = $2, title = $3, description = $4, 
		cover_url = $5, content = $6, type = $7, position = $8, duration = $9, is_free = $10, 
		status = $11, quiz_settings = $12, updated_at = NOW() WHERE id = $13 AND deleted_at IS NULL`

	cmtTag, err := q.Exec(ctx, query, entity.CourseID, entity.SectionID, entity.Title, entity.Description,
		entity.CoverURL, entity.Content, entity.Type, entity.Position, entity.Duration, entity.IsFree, status, entity.QuizSettings, entity.ID)

	if err != nil {
		return fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	if cmtTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, errorsAPP.ErrNothingToUpdate)
	}

	return nil
}
