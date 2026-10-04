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

func (s *lessonService) UpdateLesson(ctx context.Context, userID int64, lesson *domain.Lesson) error {
	const op = "service.lesson.UpdateLesson"

	existingLesson, err := s.lessonrepo.GetLessonByID(ctx, lesson.ID)
	if err != nil {
		if errors.Is(err, errorsAPP.ErrNotFoundLesson) {
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
		slog.Warn("unauthorized lesson update attempt",
			"userID", userID,
			"courseID", existingLesson.CourseID,
			"lessonID", lesson.ID,
			"action", "UpdateLesson",
		)
		return fmt.Errorf("%s: %w", op, errorsAPP.ErrForbidden)
	}

	// MERGE fields for partial update
	if lesson.Title == "" {
		lesson.Title = existingLesson.Title
	}
	if lesson.Description == "" {
		lesson.Description = existingLesson.Description
	}
	if lesson.CourseID == 0 {
		lesson.CourseID = existingLesson.CourseID
	}
	if lesson.SectionID == nil {
		lesson.SectionID = existingLesson.SectionID
	}
	if lesson.Type == "" {
		lesson.Type = existingLesson.Type
	}
	if lesson.CoverURL == "" {
		lesson.CoverURL = existingLesson.CoverURL
	}
	if lesson.Content == "" {
		lesson.Content = existingLesson.Content
	}
	
	if err := utils.ValidateLesson(int64(lesson.CourseID), lesson.Title, lesson.Description); err != nil {
		return fmt.Errorf("%s: validation failed: %w", op, err)
	}

	err = s.lessonrepo.UpdateLesson(ctx, lesson)
	if err != nil {
		if errors.Is(err, errorsAPP.ErrNothingToUpdate) {
			return fmt.Errorf("%s: %w", op, err)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
