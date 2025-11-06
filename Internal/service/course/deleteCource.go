package course

import (
	"context"
	errorsAPP "edtech/pkg/errors"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (s *service) DeleteCourse(ctx context.Context, courseID int64, userID int64) error {

	course, err := s.courserepo.GetCourseByID(ctx, courseID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errorsAPP.ErrNotFoundCourse
		}
		return fmt.Errorf("%w: %v", errorsAPP.ErrInternalDB, err)
	}

	canDelete, err := s.accessService.CanDeleteCourse(ctx, course, userID)
	if err != nil {
		return fmt.Errorf("%v: %w", errorsAPP.ErrCheckingPermissions, err)
	}

	if !canDelete {
		return errorsAPP.ErrForbidden
	}

	if err := s.courserepo.DeleteCourse(ctx, courseID); err != nil {
		return fmt.Errorf("%w: %v", errorsAPP.ErrInternalDB, err)
	}

	return nil

}
