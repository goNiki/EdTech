package lesson

import (
	"context"
	"edtech/internal/infrastructure/db"
	errorsAPP "edtech/pkg/errors"
	"fmt"
)

func (r *repository) DeleteLessonByID(ctx context.Context, q db.QueryExecutor, lessonID int64) error {
	const op = "repositiry.couse.lessonrepo.DeleteLessonByID"

	query := `DELETE FROM lessons WHERE id = $1`

	tag, err := q.Exec(ctx, query, lessonID)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, errorsAPP.ErrNotFoundLesson)
	}

	return nil
}
