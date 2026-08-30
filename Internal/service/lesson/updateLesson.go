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

	if _, err := s.lessonrepo.GetLessonByID(ctx, s.db, lesson.ID); err != nil {
		if errors.Is(err, errorsAPP.ErrNotFoundLesson) {
			return fmt.Errorf("%s: %w", op, err)
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	
	if err := utils.ValidateLesson(lesson.ID, lesson.Title, lesson.Description); err != nil {
		return fmt.Errorf("%s: validation failed: %w", op, err)
	}

	err := s.lessonrepo.UpdateLesson(ctx, s.db, lesson)
	if err != nil {
		if errors.Is(err, errorsAPP.ErrNothingToUpdate) {
			return fmt.Errorf("%s: %w", op, err)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
