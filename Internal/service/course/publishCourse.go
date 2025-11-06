package course

import (
	"context"
	"edtech/internal/domain"
	errorsAPP "edtech/pkg/errors"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (s *service) PublishCourse(ctx context.Context, userID int64, courseID int64) error {

	course, err := s.courserepo.GetCourseByID(ctx, courseID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errorsAPP.ErrNotFoundCourse
		}
		return fmt.Errorf("%w: %v", errorsAPP.ErrInternalDB, err)
	}

	canPublish, err := s.accessService.CanAccessCourseObject(ctx, course, userID, domain.ActionPublish)
	if err != nil {
		return fmt.Errorf("%v: %w", errorsAPP.ErrCheckingPermissions, err)
	}

	if !canPublish {
		return errorsAPP.ErrForbidden
	}

	//TODO еще дополнительную проверку, чтобы был хотя бы 1 урок на курсе, чтобы его можно было опубликовать
	err = course.Publish()
	if err != nil {
		return err
	}

	if err := s.courserepo.PublishCourse(ctx, course); err != nil {
		return fmt.Errorf("%w: %v", errorsAPP.ErrInternalDB, err)
	}

	return nil

}
