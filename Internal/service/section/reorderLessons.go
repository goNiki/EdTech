package section

import (
	"context"
	"fmt"

	"edtech/internal/infrastructure/db"

	"github.com/jackc/pgx/v5"
)

func (s *sectionService) ReorderLessons(ctx context.Context, sectionID int64, lessonIDs []int64) error {
	const op = "service.section.ReorderLessons"

	err := s.txManager.WithTX(ctx, pgx.TxOptions{}, func(ctx context.Context, q db.QueryExecutor) error {
		return s.lessonRepo.ReorderLessons(ctx, q, &sectionID, lessonIDs)
	})

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
