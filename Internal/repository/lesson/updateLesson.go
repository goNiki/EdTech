package lesson

import (
	"context"
	"fmt"

	"edtech/internal/domain"
	"edtech/internal/infrastructure/txmanager"
	errorsAPP "edtech/pkg/errors"
)

func (r *repository) UpdateLesson(ctx context.Context, lesson *domain.Lesson) error {
	const op = "repository.lesson.update"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	status := lesson.Status
	if status == "" {
		status = domain.StatusDraft
	}

	query := `UPDATE lessons SET course_id = $1, section_id = $2, title = $3, description = $4, 
		cover_url = $5, content = $6, type = $7, position = $8, duration = $9, is_free = $10, 
		status = $11, updated_at = NOW() WHERE id = $12 AND deleted_at IS NULL`

	cmtTag, err := q.Exec(ctx, query, lesson.CourseID, lesson.SectionID, lesson.Title, lesson.Description,
		lesson.CoverURL, lesson.Content, lesson.Type, lesson.Position, lesson.Duration, lesson.IsFree, status, lesson.ID)

	if err != nil {
		return fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	if cmtTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, errorsAPP.ErrNothingToUpdate)
	}

	return nil
}
