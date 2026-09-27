package course

import (
	"context"
	"edtech/internal/domain"
	"edtech/internal/infrastructure/db"
	errorsAPP "edtech/pkg/errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (s *service) UpdateCourseStatus(ctx context.Context, userID int64, courseID int64, status string) error {
	const op = "service.course.UpdateCourseStatus"

	if status == domain.StatusPublished {
		return s.PublishCourse(ctx, userID, courseID)
	}

	course, err := s.courserepo.GetCourseByID(ctx, s.db, courseID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	canEdit, err := s.accessService.CanEditCourse(ctx, course, userID)
	if err != nil {
		return fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrCheckingPermissions, err)
	}
	if !canEdit {
		return fmt.Errorf("%s: %w", op, errorsAPP.ErrForbidden)
	}

	err = s.txManager.WithTX(ctx, pgx.TxOptions{}, func(ctx context.Context, q db.QueryExecutor) error {
		if err := s.courserepo.UpdateCourseStatus(ctx, q, courseID, status); err != nil {
			return err
		}
		if err := s.sectionrepo.UpdateStatusByCourseID(ctx, q, courseID, status); err != nil {
			return err
		}
		return s.lessonrepo.UpdateStatusByCourseID(ctx, q, courseID, status)
	})

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
