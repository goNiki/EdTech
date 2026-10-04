package enrollment

import (
	"context"
	"errors"
	"fmt"

	errorsAPP "edtech/pkg/errors"

	"github.com/jackc/pgx/v5"
)

func (s *service) UnenrollUser(ctx context.Context, userID int64, courseID int64) error {
	const op = "service.enrollment.UnenrollUser"

	role, err := s.enrolledrepo.GetRoleUserInCourse(ctx, userID, courseID)
	if err != nil {
		if errors.Is(err, errorsAPP.ErrNotFoundCourse) {
			return fmt.Errorf("%s: %w", op, errorsAPP.ErrNotEnrolled)
		}
		return fmt.Errorf("%s: get role: %w", op, err)
	}

	if role == "" {
		return fmt.Errorf("%s: %w", op, errorsAPP.ErrNotEnrolled)
	}
	if role == "creator" {
		return fmt.Errorf("%s: %w", op, errorsAPP.ErrCreatorCannotUnenroll)
	}

	err = s.txManager.WithTX(ctx, pgx.TxOptions{}, func(ctx context.Context) error {
		if err := s.enrolledrepo.UnenrollUser(ctx, userID, courseID); err != nil {
			return err
		}
		if err := s.courserepo.DecrementEnrolledCount(ctx, courseID); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *service) TeacherUnenrollUser(ctx context.Context, teacherID int64, targetUserID int64, courseID int64) error {
	const op = "service.enrollment.TeacherUnenrollUser"

	course, err := s.courserepo.GetCourseByID(ctx, courseID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	canManage, err := s.accessService.CanManageCourseUsers(ctx, course, teacherID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if !canManage {
		return fmt.Errorf("%s: %w", op, errorsAPP.ErrForbidden)
	}

	return s.UnenrollUser(ctx, targetUserID, courseID)
}

