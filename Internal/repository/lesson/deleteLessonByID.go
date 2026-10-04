package lesson

import (
	"context"
	"edtech/internal/infrastructure/txmanager"
	errorsAPP "edtech/pkg/errors"
	"fmt"
)

func (r *repository) DeleteLessonByID(ctx context.Context, lessonID int64) error {
	const op = "repositiry.couse.lessonrepo.DeleteLessonByID"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	query := `DELETE FROM lessons WHERE id = $1`

	tag, err := q.Exec(ctx, query, lessonID)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, errorsAPP.ErrLessonNotFound)
	}

	return nil
}
