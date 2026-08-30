package lesson

import (
	"context"
	errorsAPP "edtech/pkg/errors"
	"errors"
	"fmt"
)

func (s *service) DeleteLesson(ctx context.Context, lessonID int64) error {

	const op = "service.lesson.DeleteLesson"

	if err := s.lessonrepo.DeleteLessonByID(ctx, s.db, lessonID); err != nil {
		if errors.Is(err, errorsAPP.ErrNotFoundLesson) {
			return fmt.Errorf("%s: %w", op, err)
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
