package course

import (
	"context"
	"edtech/internal/domain"
	errorsAPP "edtech/pkg/errors"
	"edtech/pkg/utils"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (s *service) UpdateCourse(ctx context.Context, course *domain.Course, userID int64) error {

	if err := utils.ValidateCourse(course.Title, course.Slug, course.CreatedBy, course.Visibility, course.Status); err != nil {
		return err
	}

	canEdit, err := s.accessService.CanEditCourse(ctx, course, userID)
	if err != nil {
		return fmt.Errorf("%v: %w", errorsAPP.ErrCheckingPermissions, err)
	}

	if !canEdit {
		return errorsAPP.ErrForbidden
	}

	oldcourse, err := s.courserepo.GetCourseByID(ctx, course.Id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errorsAPP.ErrNotFoundCourse
		}
		return fmt.Errorf("%w: %v", errorsAPP.ErrInternalDB, err)
	}

	if course.Slug != oldcourse.Slug {
		if _, err := s.courserepo.GetCourseBySlug(ctx, course.Slug); !errors.Is(err, pgx.ErrNoRows) {
			if err == nil {
				return errorsAPP.ErrSlugAlreadyExists
			}
			return fmt.Errorf("%w: %v", errorsAPP.ErrInternalDB, err)
		}
	}

	if err := s.courserepo.UpdateCourse(ctx, course); err != nil {
		if errors.Is(err, errorsAPP.ErrNothingToUpdate) {
			return errorsAPP.ErrNothingToUpdate
		}
		return fmt.Errorf("%w: %v", errorsAPP.ErrInternalDB, err)
	}

	return nil

}
