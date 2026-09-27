package section

import (
	"context"
	"fmt"

	"edtech/internal/domain"
	"edtech/internal/infrastructure/db"

	"github.com/jackc/pgx/v5"
)

func (s *sectionService) CreateSection(ctx context.Context, section *domain.Section) (*domain.Section, error) {
	const op = "service.section.CreateSection"

	var createdSection *domain.Section

	err := s.txManager.WithTX(ctx, pgx.TxOptions{}, func(ctx context.Context, q db.QueryExecutor) error {
		if section.Position == 0 {
			maxPos, err := s.sectionRepo.GetMaxPositionByCourseID(ctx, q, section.CourseID)
			if err != nil {
				return err
			}
			section.Position = maxPos + 1
		}
		if section.Status == "" {
			section.Status = domain.StatusDraft
		}

		var err error
		createdSection, err = s.sectionRepo.CreateSection(ctx, q, section)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return createdSection, nil
}
