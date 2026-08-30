package progress

import (
	"context"
	"fmt"

	"edtech/internal/domain"
	"edtech/internal/infrastructure/db"

	"github.com/jackc/pgx/v5"
)

func (s *service) UpdateLessonProgress(ctx context.Context, userID int64, lessonID int64, input domain.UpdateProgressInput) error {
	const op = "service.progress.UpdateLessonProgress"

	err := s.txManager.WithTX(ctx, pgx.TxOptions{}, func(ctx context.Context, tx db.QueryExecutor) error {
		if input.TimeSpent > 0 || input.LastPosition > 0 {
			if err := s.progressRepo.UpdateLessonProgressTime(ctx, tx, userID, lessonID, input.TimeSpent, input.LastPosition); err != nil {
				return err
			}
		}
		if input.Status != "" {
			if err := s.progressRepo.UpdateLessonProgressStatus(ctx, tx, userID, lessonID, input.Status); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
