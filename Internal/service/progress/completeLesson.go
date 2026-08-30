package progress

import (
	"context"
	"fmt"

	"edtech/internal/domain"
)

func (s *service) CompleteLesson(ctx context.Context, userID int64, lessonID int64) error {
	const op = "service.progress.CompleteLesson"

	if err := s.progressRepo.UpdateLessonProgressStatus(ctx, s.db, userID, lessonID, domain.ProgressStatusCompleted); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
