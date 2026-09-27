package lesson

import (
	"context"
	"fmt"

	"edtech/internal/infrastructure/db"
	errorsAPP "edtech/pkg/errors"
)

func (r *repository) ReorderLessons(ctx context.Context, q db.QueryExecutor, sectionID *int64, lessonIDs []int64) error {
	const op = "repository.lesson.ReorderLessons"

	for idx, id := range lessonIDs {
		pos := int64(idx + 1)
		var query string
		var args []interface{}

		if sectionID != nil {
			query = `UPDATE lessons SET position = $1, section_id = $2, updated_at = NOW() WHERE id = $3 AND deleted_at IS NULL`
			args = []interface{}{pos, *sectionID, id}
		} else {
			query = `UPDATE lessons SET position = $1, updated_at = NOW() WHERE id = $2 AND deleted_at IS NULL`
			args = []interface{}{pos, id}
		}

		_, err := q.Exec(ctx, query, args...)
		if err != nil {
			return fmt.Errorf("%s: update position for lesson %d: %w: %w", op, id, errorsAPP.ErrInternalDB, err)
		}
	}

	return nil
}
