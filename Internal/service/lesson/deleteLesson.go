package lesson

import (
	"context"
	errorsAPP "edtech/pkg/errors"
	"errors"
	"fmt"
	"log/slog"
)

func (s *lessonService) DeleteLesson(ctx context.Context, userID int64, lessonID int64) error {

	const op = "service.lesson.DeleteLesson"

	// Load the lesson to find its courseID
	existingLesson, err := s.lessonrepo.GetLessonByID(ctx, lessonID)
	if err != nil {
		if errors.Is(err, errorsAPP.ErrLessonNotFound) {
			return fmt.Errorf("%s: %w", op, err)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	// RBAC: check that user can edit the course this lesson belongs to
	course, err := s.courserepo.GetCourseByID(ctx, int64(existingLesson.CourseID))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	canEdit, err := s.accessService.CanEditCourse(ctx, course, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if !canEdit {
		slog.Warn("unauthorized lesson deletion attempt",
			"userID", userID,
			"courseID", existingLesson.CourseID,
			"lessonID", lessonID,
			"action", "DeleteLesson",
		)
		return fmt.Errorf("%s: %w", op, errorsAPP.ErrForbidden)
	}

	if err := s.lessonrepo.DeleteLessonByID(ctx, lessonID); err != nil {
		if errors.Is(err, errorsAPP.ErrLessonNotFound) {
			return fmt.Errorf("%s: %w", op, err)
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
