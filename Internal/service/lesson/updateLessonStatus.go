package lesson

import (
	"context"
	errorsAPP "edtech/pkg/errors"
	"errors"
	"fmt"
	"log/slog"
)

func (s *lessonService) UpdateLessonStatus(ctx context.Context, userID int64, lessonID int64, status string) error {
	const op = "service.lesson.UpdateLessonStatus"

	// Load the lesson to find its courseID
	existingLesson, err := s.lessonrepo.GetLessonByID(ctx, s.db, lessonID)
	if err != nil {
		if errors.Is(err, errorsAPP.ErrNotFoundLesson) {
			return fmt.Errorf("%s: %w", op, err)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	// RBAC: check that user can edit the course this lesson belongs to
	course, err := s.courserepo.GetCourseByID(ctx, s.db, int64(existingLesson.CourseID))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	canEdit, err := s.accessService.CanEditCourse(ctx, course, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if !canEdit {
		slog.Warn("unauthorized lesson status update attempt",
			"userID", userID,
			"courseID", existingLesson.CourseID,
			"lessonID", lessonID,
			"action", "UpdateLessonStatus",
		)
		return fmt.Errorf("%s: %w", op, errorsAPP.ErrForbidden)
	}

	if err := s.lessonrepo.UpdateLessonStatus(ctx, s.db, lessonID, status); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
