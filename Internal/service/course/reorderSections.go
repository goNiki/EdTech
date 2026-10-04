package course

import (
	"context"
	"fmt"

	errorsAPP "edtech/pkg/errors"
)

func (s *service) ReorderSections(ctx context.Context, userID, courseID int64, sectionIDs []int64) error {
	const op = "service.course.ReorderSections"

	course, err := s.courserepo.GetCourseByID(ctx, courseID)
	if err != nil {
		return fmt.Errorf("%s: get course: %w", op, err)
	}

	canEdit, err := s.accessService.CanEditCourse(ctx, course, userID)
	if err != nil {
		return fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrCheckingPermissions, err)
	}
	if !canEdit {
		return fmt.Errorf("%s: %w", op, errorsAPP.ErrForbidden)
	}

	if err := s.sectionrepo.ReorderSections(ctx, courseID, sectionIDs); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
