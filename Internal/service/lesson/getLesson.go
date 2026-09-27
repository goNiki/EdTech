package lesson

import (
	"context"
	"edtech/internal/domain"
	"fmt"
)

func (s *service) GetLesson(ctx context.Context, id int64) (*domain.Lesson, error) {
	const op = "service.lesson.GetLesson"

	lesson, err := s.lessonrepo.GetLessonByID(ctx, s.db, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return lesson, nil
}
