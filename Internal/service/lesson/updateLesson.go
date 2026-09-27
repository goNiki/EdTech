package lesson

import (
	"context"
	"edtech/internal/domain"
	errorsAPP "edtech/pkg/errors"
	"edtech/pkg/utils"
	"errors"
	"fmt"
)

func (s *service) UpdateLesson(ctx context.Context, lesson *domain.Lesson) error {
	const op = "service.lesson.UpdateLesson"

	existingLesson, err := s.lessonrepo.GetLessonByID(ctx, s.db, lesson.ID)
	if err != nil {
		if errors.Is(err, errorsAPP.ErrNotFoundLesson) {
			return fmt.Errorf("%s: %w", op, err)
		}
		return fmt.Errorf("%s: %w", op, err)
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

	err = s.lessonrepo.UpdateLesson(ctx, s.db, lesson)
	if err != nil {
		if errors.Is(err, errorsAPP.ErrNothingToUpdate) {
			return fmt.Errorf("%s: %w", op, err)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
