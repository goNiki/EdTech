package section

import (
	"context"
	"fmt"
	"log/slog"

	"edtech/internal/infrastructure/db"
	errorsAPP "edtech/pkg/errors"

	"github.com/jackc/pgx/v5"
)

func (s *sectionService) DeleteSection(ctx context.Context, userID int64, sectionID int64) error {
	const op = "service.section.DeleteSection"

	err := s.txManager.WithTX(ctx, pgx.TxOptions{}, func(ctx context.Context, q db.QueryExecutor) error {
		// Load the section to find its courseID
		existing, err := s.sectionRepo.GetSectionByID(ctx, q, sectionID)
		if err != nil {
			return err
		}

		// RBAC: check that user can edit the course
		course, err := s.courseRepo.GetCourseByID(ctx, q, existing.CourseID)
		if err != nil {
			return err
		}

		canEdit, err := s.accessService.CanEditCourse(ctx, course, userID)
		if err != nil {
			return err
		}
		if !canEdit {
			slog.Warn("unauthorized section deletion attempt",
				"userID", userID,
				"courseID", existing.CourseID,
				"sectionID", sectionID,
				"action", "DeleteSection",
			)
			return errorsAPP.ErrForbidden
		}

		return s.sectionRepo.DeleteSection(ctx, q, sectionID)
	})

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
