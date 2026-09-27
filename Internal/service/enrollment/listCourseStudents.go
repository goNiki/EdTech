package enrollment

import (
	"context"
	"fmt"

	"edtech/internal/domain"
	errorsAPP "edtech/pkg/errors"
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

func (s *service) ListCourseStudentsWithProgress(ctx context.Context, teacherID, courseID int64, page, pageSize int64) ([]domain.CourseStudentItem, int64, error) {
	const op = "service.enrollment.ListCourseStudentsWithProgress"

	course, err := s.courserepo.GetCourseByID(ctx, s.db, courseID)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: get course: %w", op, err)
	}

	canManage, err := s.accessService.CanManageCourseUsers(ctx, course, teacherID)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrCheckingPermissions, err)
	}
	if !canManage {
		return nil, 0, fmt.Errorf("%s: %w", op, errorsAPP.ErrForbidden)
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	students, total, err := s.enrolledrepo.ListCourseStudentsWithProgress(ctx, s.db, courseID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	return students, total, nil
}
