package section

import (
	"context"
	"edtech/internal/domain"
	"edtech/internal/infrastructure/db"
	"edtech/internal/infrastructure/txmanager"
	"edtech/internal/repository"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type sectionService struct {
	sectionRepo repository.SectionRepository
	lessonRepo  repository.LessonRepository
	txManager   txmanager.TransactionManager
}

func NewSectionService(
	sectionRepo repository.SectionRepository,
	lessonRepo repository.LessonRepository,
	txManager txmanager.TransactionManager,
) *sectionService {
	return &sectionService{
		sectionRepo: sectionRepo,
		lessonRepo:  lessonRepo,
		txManager:   txManager,
	}
}

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
