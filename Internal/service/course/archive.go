package course

import (
	"context"
	"fmt"
	"time"

	errorsAPP "edtech/pkg/errors"
)

func (s *service) ArchiveCourse(ctx context.Context, userID int64, courseID int64) error {
	const op = "service.course.ArchiveCourse"

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

	if err := course.CanArchive(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	course.Archive(time.Now())

	if err := s.courserepo.ArchiveCourse(ctx, course); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
