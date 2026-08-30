package lesson

import (
	"context"
	"edtech/internal/infrastructure/db"
	"fmt"
)

func (r *repository) GetMaxPositionByCourseID(ctx context.Context, q db.QueryExecutor, courseId int64) (int64, error) {
	const op = "repositiry.couse.lessonrepo.getmaxpositionbycourseid"

	query := `SELECT COALESCE(MAX(position), 0) FROM lessons WHERE course_id = $1`

	var maxPosition int64

	err := q.QueryRow(ctx, query, courseId).Scan(&maxPosition)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return maxPosition, nil
}
