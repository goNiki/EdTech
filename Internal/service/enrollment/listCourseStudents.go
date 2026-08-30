package enrollment

import (
	"context"
	"fmt"

	"edtech/internal/domain"
)

func (s *service) ListCourseStudents(ctx context.Context, courseID int64, page int64, pageSize int64) ([]domain.User, int, error) {
	const op = "service.enrollment.ListCourseStudents"

	if _, err := s.courserepo.GetCourseByID(ctx, s.db, courseID); err != nil {
		return nil, 0, fmt.Errorf("%s: get course: %w", op, err)
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	students, total, err := s.enrolledrepo.ListCourseStudents(ctx, s.db, courseID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	return students, total, nil
}
