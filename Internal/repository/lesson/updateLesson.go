package lesson

import (
	"context"
	"edtech/internal/domain"
	errorsAPP "edtech/pkg/errors"
	"fmt"
)

// Обновляет данные урока, если ничего не обновил, то выдает ошибку. чтобы в дальнейшем обработать в хенделере.
func (r *repository) UpdateLesson(ctx context.Context, lesson *domain.Lesson) error {
	const op = "repository.course.lesson.updatelesson"

	query := `
		UPDATE lessons
		SET title = $1,
		    description = $2,
		    cover_url = $3,
		    content = $4,
		    update_at = NOW()
		WHERE id = $5
		  AND (
		      title IS DISTINCT FROM $1 OR
		      description IS DISTINCT FROM $2 OR
		      cover_url IS DISTINCT FROM $3 OR
		      content IS DISTINCT FROM $4
		  )`

	cmdTag, err := r.Pool.Exec(ctx, query, lesson.Title, lesson.Description, lesson.CoverURL, lesson.Content, lesson.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, errorsAPP.ErrNothingToUpdate)
	}

	return nil
}
