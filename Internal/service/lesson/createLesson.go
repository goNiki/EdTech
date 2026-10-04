package lesson

import (
	"context"
	"edtech/internal/domain"
	errorsAPP "edtech/pkg/errors"
	"edtech/pkg/utils"
	"errors"
	"fmt"
	"log/slog"
)

func (s *lessonService) CreateLesson(ctx context.Context, userID int64, lesson *domain.Lesson) (int64, error) {

	const op = "service.lesson.CreateLesson"

	if err := utils.ValidateLesson(int64(lesson.CourseID), lesson.Title, lesson.Description); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	course, err := s.courserepo.GetCourseByID(ctx, int64(lesson.CourseID))
	if err != nil {
		if errors.Is(err, errorsAPP.ErrNotFoundCourse) {
			return 0, fmt.Errorf("%s: %w", op, err)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	canEdit, err := s.accessService.CanEditCourse(ctx, course, userID)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	if !canEdit {
		slog.Warn("unauthorized lesson creation attempt",
			"userID", userID,
			"courseID", lesson.CourseID,
			"action", "CreateLesson",
		)
		return 0, fmt.Errorf("%s: %w", op, errorsAPP.ErrForbidden)
	}

	position, err := s.lessonrepo.GetMaxPositionByCourseID(ctx, int64(lesson.CourseID))
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	lesson.Position = position + 1

	err = s.lessonrepo.CreateLesson(ctx, lesson)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return int64(lesson.ID), nil
}
