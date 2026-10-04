package enrollment

import (
	"context"
	"errors"
	"fmt"

	errorsAPP "edtech/pkg/errors"
)

func (s *service) ChangeUserRole(ctx context.Context, courseID int64, targetUserID int64, newRole string) error {
	const op = "service.enrollment.ChangeUserRole"

	if newRole != "student" && newRole != "teacher" {
		return fmt.Errorf("%s: %w", op, errorsAPP.ErrInvalidAction)
	}

	role, err := s.enrolledrepo.GetRoleUserInCourse(ctx, targetUserID, courseID)
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

	if err := s.enrolledrepo.ChangeUserRole(ctx, courseID, targetUserID, newRole); err != nil {
		return fmt.Errorf("%s: change role: %w", op, err)
	}

	return nil
}
