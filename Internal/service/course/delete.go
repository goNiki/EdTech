package course

import (
	"context"
	"fmt"

	errorsAPP "edtech/pkg/errors"
)

func (s *service) DeleteCourse(ctx context.Context, courseID int64, userID int64) error {
	const op = "service.course.DeleteCourse"

	course, err := s.courserepo.GetCourseByID(ctx, courseID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	canDelete, err := s.accessService.CanDeleteCourse(ctx, course, userID)
	if err != nil {
		return fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrCheckingPermissions, err)
	}
	if !canDelete {
		return fmt.Errorf("%s: %w", op, errorsAPP.ErrForbidden)
	}

	if err := s.courserepo.DeleteCourse(ctx, courseID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
