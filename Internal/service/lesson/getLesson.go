package lesson

import (
	"context"
	"edtech/internal/domain"
	errorsAPP "edtech/pkg/errors"
	"errors"
	"fmt"
)

func (s *service) GetLesson(ctx context.Context, id int64) (*domain.Lesson, error) {
	const op = "service.lesson.GetLesson"

	lesson, err := s.lessonrepo.GetLessonByID(ctx, s.db, id)
	if err != nil {
		if errors.Is(err, errorsAPP.ErrNotFoundLesson) {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return lesson, nil
}
