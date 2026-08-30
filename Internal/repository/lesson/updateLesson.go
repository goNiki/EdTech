package lesson

import (
	"context"
	"edtech/internal/domain"
	"edtech/internal/infrastructure/db"
	errorsAPP "edtech/pkg/errors"
	"fmt"
)

func (r *repository) UpdateLesson(ctx context.Context, q db.QueryExecutor, lesson *domain.Lesson) error {
	const op = "repository.lesson.update"

	query := `UPDATE lessons SET course_id = $1, section_id = $2, title = $3, description = $4, 
		cover_url = $5, content = $6, type = $7, position = $8, duration = $9, is_free = $10, 
		updated_at = NOW() WHERE id = $11 AND deleted_at IS NULL`

	cmtTag, err := q.Exec(ctx, query, lesson.CourseID, lesson.SectionID, lesson.Title, lesson.Description,
		lesson.CoverURL, lesson.Content, lesson.Type, lesson.Position, lesson.Duration, lesson.IsFree, lesson.ID)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if cmtTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, errorsAPP.ErrNothingToUpdate)
	}

	return nil
}
