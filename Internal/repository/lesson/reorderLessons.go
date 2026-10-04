package lesson

import (
	"context"
	"fmt"

	"edtech/internal/infrastructure/txmanager"
	errorsAPP "edtech/pkg/errors"
)

func (r *repository) ReorderLessons(ctx context.Context, sectionID *int64, lessonIDs []int64) error {
	const op = "repository.lesson.ReorderLessons"
	if len(lessonIDs) == 0 {
		return nil
	}

	positions := make([]int32, len(lessonIDs))
	for i := range lessonIDs {
		positions[i] = int32(i + 1)
	}

	var query string
	var args []any

	if sectionID != nil {
		query = `
			UPDATE lessons AS l
			SET position = v.new_pos, section_id = $1, updated_at = NOW()
			FROM (SELECT unnest($2::bigint[]) AS id, unnest($3::int[]) AS new_pos) AS v
			WHERE l.id = v.id AND l.deleted_at IS NULL
		`
		args = []any{*sectionID, lessonIDs, positions}
	} else {
		query = `
			UPDATE lessons AS l
			SET position = v.new_pos, updated_at = NOW()
			FROM (SELECT unnest($1::bigint[]) AS id, unnest($2::int[]) AS new_pos) AS v
			WHERE l.id = v.id AND l.deleted_at IS NULL
		`
		args = []any{lessonIDs, positions}
	}

	q := txmanager.GetQueryExecutor(ctx, r.Pool)
	_, err := q.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: batch update lessons: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	return nil
}
