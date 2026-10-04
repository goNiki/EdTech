package course

import (
	"context"
	"edtech/internal/domain"
	errorsAPP "edtech/pkg/errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

func (s *service) PublishCourse(ctx context.Context, userID int64, courseID int64) error {
	const op = "service.course.PublishCourse"

	course, err := s.courserepo.GetCourseByID(ctx, courseID)
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

	course.Publish(time.Now())

	err = s.txManager.WithTX(ctx, pgx.TxOptions{}, func(ctx context.Context) error {
		if err := s.courserepo.PublishCourse(ctx, course); err != nil {
			return err
		}
		if err := s.sectionrepo.UpdateStatusByCourseID(ctx, courseID, domain.StatusPublished); err != nil {
			return err
		}
		return s.lessonrepo.UpdateStatusByCourseID(ctx, courseID, domain.StatusPublished)
	})

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
