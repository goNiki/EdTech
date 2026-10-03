package section

import (
	"context"
	"fmt"
	"log/slog"

	"edtech/internal/domain"
	"edtech/internal/infrastructure/db"
	errorsAPP "edtech/pkg/errors"

	"github.com/jackc/pgx/v5"
)

func (s *sectionService) CreateSection(ctx context.Context, userID int64, section *domain.Section) (*domain.Section, error) {
	const op = "service.section.CreateSection"

	// RBAC: check that user can edit the course
	course, err := s.courseRepo.GetCourseByID(ctx, s.db, section.CourseID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	canEdit, err := s.accessService.CanEditCourse(ctx, course, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if !canEdit {
		slog.Warn("unauthorized section creation attempt",
			"userID", userID,
			"courseID", section.CourseID,
			"action", "CreateSection",
		)
		return nil, fmt.Errorf("%s: %w", op, errorsAPP.ErrForbidden)
	}

	var createdSection *domain.Section

	err = s.txManager.WithTX(ctx, pgx.TxOptions{}, func(ctx context.Context, q db.QueryExecutor) error {
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
