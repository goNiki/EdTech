package section

import (
	"context"
	"fmt"
	"log/slog"

	errorsAPP "edtech/pkg/errors"

	"github.com/jackc/pgx/v5"
)

func (s *sectionService) ReorderLessons(ctx context.Context, userID int64, sectionID int64, lessonIDs []int64) error {
	const op = "service.section.ReorderLessons"

	err := s.txManager.WithTX(ctx, pgx.TxOptions{}, func(ctx context.Context) error {
		// Load the section to find its courseID
		existing, err := s.sectionRepo.GetSectionByID(ctx, sectionID)
		if err != nil {
			return err
		}

		// RBAC: check that user can edit the course
		course, err := s.courseRepo.GetCourseByID(ctx, existing.CourseID)
		if err != nil {
			return err
		}

		canEdit, err := s.accessService.CanEditCourse(ctx, course, userID)
		if err != nil {
			return err
		}
		if !canEdit {
			slog.Warn("unauthorized lesson reorder attempt",
				"userID", userID,
				"courseID", existing.CourseID,
				"sectionID", sectionID,
				"action", "ReorderLessons",
			)
			return errorsAPP.ErrForbidden
		}

		return s.lessonRepo.ReorderLessons(ctx, &sectionID, lessonIDs)
	})

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
