package section

import (
	"context"
	"fmt"

	"edtech/internal/domain"
	"edtech/internal/infrastructure/db"

	"github.com/jackc/pgx/v5"
)

func (s *sectionService) UpdateSection(ctx context.Context, section *domain.Section) error {
	const op = "service.section.UpdateSection"

	err := s.txManager.WithTX(ctx, pgx.TxOptions{}, func(ctx context.Context, q db.QueryExecutor) error {
		existing, err := s.sectionRepo.GetSectionByID(ctx, q, section.ID)
		if err != nil {
			return err
		}

		if section.Title != "" {
			existing.Title = section.Title
		}
		if section.Description != "" {
			existing.Description = section.Description
		}
		if section.Position != 0 {
			existing.Position = section.Position
		}
		if section.Status != "" {
			existing.Status = section.Status
		}

		return s.sectionRepo.UpdateSection(ctx, q, existing)
	})

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
