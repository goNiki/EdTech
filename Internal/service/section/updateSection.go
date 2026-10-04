package section

import (
	"context"
	"fmt"
	"log/slog"

	"edtech/internal/domain"
	errorsAPP "edtech/pkg/errors"

	"github.com/jackc/pgx/v5"
)

func (s *sectionService) UpdateSection(ctx context.Context, userID int64, section *domain.Section) error {
	const op = "service.section.UpdateSection"

	err := s.txManager.WithTX(ctx, pgx.TxOptions{}, func(ctx context.Context) error {
		existing, err := s.sectionRepo.GetSectionByID(ctx, section.ID)
		if err != nil {
			return err
		}

		// RBAC: check that user can edit the course this section belongs to
		course, err := s.courseRepo.GetCourseByID(ctx, existing.CourseID)
		if err != nil {
			return err
		}

		canEdit, err := s.accessService.CanEditCourse(ctx, course, userID)
		if err != nil {
			return err
		}
		if !canEdit {
			slog.Warn("unauthorized section update attempt",
				"userID", userID,
				"courseID", existing.CourseID,
				"sectionID", section.ID,
				"action", "UpdateSection",
			)
			return errorsAPP.ErrForbidden
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

		return s.sectionRepo.UpdateSection(ctx, existing)
	})

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
