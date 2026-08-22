package auth

import (
	"context"
	"fmt"

	"edtech/internal/domain"
	"edtech/internal/infrastructure/db"
	errorsAPP "edtech/pkg/errors"

	"github.com/jackc/pgx/v5"
)

func (s *service) ChangeUserRole(ctx context.Context, adminID int64, targetUserID int64, newRole domain.Role) error {
	const op = "service.auth.ChangeUserRole"

	admin, err := s.repo.GetUserByID(ctx, s.db, adminID)
	if err != nil {
		return fmt.Errorf("%s: get admin: %w", op, err)
	}
	if !admin.IsAdmin() {
		return errorsAPP.ErrForbidden
	}

	if !newRole.IsValid() {
		return errorsAPP.ErrInvalidRole
	}

	if adminID == targetUserID && newRole != domain.RoleAdmin {
		return errorsAPP.ErrCannotModifySelf
	}

	if _, err := s.repo.GetUserByID(ctx, s.db, targetUserID); err != nil {
		return fmt.Errorf("%s: get target user: %w", op, err)
	}

	if err := s.repo.UpdateRole(ctx, s.db, targetUserID, newRole); err != nil {
		return fmt.Errorf("%s: update role: %w", op, err)
	}

	return nil
}

func (s *service) SetUserBanned(ctx context.Context, adminID int64, targetUserID int64, isBanned bool) error {
	const op = "service.auth.SetUserBanned"

	admin, err := s.repo.GetUserByID(ctx, s.db, adminID)
	if err != nil {
		return fmt.Errorf("%s: get admin: %w", op, err)
	}
	if !admin.IsAdmin() {
		return errorsAPP.ErrForbidden
	}

	if adminID == targetUserID {
		return errorsAPP.ErrCannotModifySelf
	}

	if _, err := s.repo.GetUserByID(ctx, s.db, targetUserID); err != nil {
		return fmt.Errorf("%s: get target user: %w", op, err)
	}

	err = s.txmanager.WithTX(ctx, pgx.TxOptions{}, func(ctx context.Context, q db.QueryExecutor) error {
		if err := s.repo.SetBannedStatus(ctx, q, targetUserID, isBanned); err != nil {
			return err
		}

		if isBanned {
			if err := s.refreshRepo.DeleteAllByUserID(ctx, q, targetUserID); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
