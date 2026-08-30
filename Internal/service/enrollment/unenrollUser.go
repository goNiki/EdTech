package enrollment

import (
	"context"
	"errors"
	"fmt"

	"edtech/internal/infrastructure/db"
	errorsAPP "edtech/pkg/errors"

	"github.com/jackc/pgx/v5"
)

func (s *service) UnenrollUser(ctx context.Context, userID int64, courseID int64) error {
	const op = "service.enrollment.UnenrollUser"

	role, err := s.enrolledrepo.GetRoleUserInCourse(ctx, s.db, userID, courseID)
	if err != nil {
		if errors.Is(err, errorsAPP.ErrNotFoundCourse) {
			return fmt.Errorf("%s: %w", op, errorsAPP.ErrNotEnrolled)
		}
		return fmt.Errorf("%s: get role: %w", op, err)
	}

	if role == "" {
		return fmt.Errorf("%s: %w", op, errorsAPP.ErrNotEnrolled)
	}
	if role == "creator" {
		return fmt.Errorf("%s: %w", op, errorsAPP.ErrCreatorCannotUnenroll)
	}

	err = s.txManager.WithTX(ctx, pgx.TxOptions{}, func(ctx context.Context, tx db.QueryExecutor) error {
		if err := s.enrolledrepo.UnenrollUser(ctx, tx, userID, courseID); err != nil {
			return err
		}
		if err := s.courserepo.DecrementEnrolledCount(ctx, tx, courseID); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
