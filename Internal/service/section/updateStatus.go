package section

import (
	"context"
	"fmt"

	"edtech/internal/infrastructure/db"

	"github.com/jackc/pgx/v5"
)

func (s *sectionService) UpdateSectionStatus(ctx context.Context, sectionID int64, status string) error {
	const op = "service.section.UpdateSectionStatus"

	err := s.txManager.WithTX(ctx, pgx.TxOptions{}, func(ctx context.Context, q db.QueryExecutor) error {
		if err := s.sectionRepo.UpdateSectionStatus(ctx, q, sectionID, status); err != nil {
			return err
		}
		return s.lessonRepo.UpdateStatusBySectionID(ctx, q, sectionID, status)
	})

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
