package course

import (
	"context"
	"fmt"

	"edtech/internal/domain"
	"edtech/internal/infrastructure/txmanager"
	errorsAPP "edtech/pkg/errors"
)

func (r *repository) CountCourses(ctx context.Context, filter domain.CourseFilter) (int64, error) {
	const op = "repository.course.CountCourses"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	whereClause, args := buildCourseFilterQuery(filter)
	query := "SELECT COUNT(*) FROM courses " + whereClause

	var total int64
	err := q.QueryRow(ctx, query, args...).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	return total, nil
}
