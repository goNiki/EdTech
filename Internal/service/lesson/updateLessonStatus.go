package lesson

import (
	"context"
	"fmt"
)

func (s *service) UpdateLessonStatus(ctx context.Context, lessonID int64, status string) error {
	const op = "service.lesson.UpdateLessonStatus"

	if err := s.lessonrepo.UpdateLessonStatus(ctx, s.db, lessonID, status); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
