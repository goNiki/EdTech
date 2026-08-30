package course

import (
	"context"
	"fmt"

	"edtech/internal/domain"
	errorsAPP "edtech/pkg/errors"
)

func (s *service) UpdateCourse(ctx context.Context, courseID int64, userID int64, input domain.UpdateCourseInput) (*domain.Course, error) {
	const op = "service.course.UpdateCourse"

	course, err := s.courserepo.GetCourseByID(ctx, s.db, courseID)
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

	course.Update(input)

	if err := course.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := s.courserepo.UpdateCourse(ctx, s.db, course); err != nil {
		return nil, fmt.Errorf("%s: save updated course: %w", op, err)
	}

	return course, nil
}
