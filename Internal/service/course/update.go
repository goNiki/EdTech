package course

import (
	"context"
	"fmt"

	"edtech/internal/domain"
	errorsAPP "edtech/pkg/errors"
	"edtech/pkg/utils"
)

func (s *service) UpdateCourse(ctx context.Context, courseID int64, userID int64, input domain.UpdateCourseInput) (*domain.Course, error) {
	const op = "service.course.UpdateCourse"

	course, err := s.courserepo.GetCourseByID(ctx, courseID)
	if err != nil {
		return nil, fmt.Errorf("%s: get course: %w", op, err)
	}

	canEdit, err := s.accessService.CanEditCourse(ctx, course, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrCheckingPermissions, err)
	}
	if !canEdit {
		return nil, fmt.Errorf("%s: %w", op, errorsAPP.ErrForbidden)
	}

	if input.Slug != nil {
		cleanSlug := utils.NormalizeSlug(*input.Slug)
		if cleanSlug != course.Slug {
			uniqueSlug, err := s.resolveUniqueSlug(ctx, cleanSlug)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", op, err)
			}
			input.Slug = &uniqueSlug
		} else {
			input.Slug = &cleanSlug
		}
	}

	course.Update(input)

	if err := course.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := s.courserepo.UpdateCourse(ctx, course); err != nil {
		return nil, fmt.Errorf("%s: save updated course: %w", op, err)
	}

	return course, nil
}
