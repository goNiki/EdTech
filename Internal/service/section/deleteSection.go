package section

import (
	"context"
	"fmt"

	"edtech/internal/infrastructure/db"

	"github.com/jackc/pgx/v5"
)

func (s *sectionService) DeleteSection(ctx context.Context, sectionID int64) error {
	const op = "service.section.DeleteSection"

	err := s.txManager.WithTX(ctx, pgx.TxOptions{}, func(ctx context.Context, q db.QueryExecutor) error {
		return s.sectionRepo.DeleteSection(ctx, q, sectionID)
	})

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
