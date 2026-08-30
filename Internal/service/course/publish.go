package course

import (
	"context"
	errorsAPP "edtech/pkg/errors"
	"fmt"
	"time"
)

func (s *service) PublishCourse(ctx context.Context, userID int64, courseID int64) error {
	const op = "service.course.PublishCourse"

	course, err := s.courserepo.GetCourseByID(ctx, s.db, courseID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	canPublish, err := s.accessService.CanPublishCourse(ctx, course, userID)
	if err != nil {
		return fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrCheckingPermissions, err)
	}

	if !canPublish {
		return fmt.Errorf("%s: %w", op, errorsAPP.ErrForbidden)
	}

	if err := course.CanPublish(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	publishTime := time.Now

	course.Publish(publishTime())

	if err := s.courserepo.PublishCourse(ctx, s.db, course); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
