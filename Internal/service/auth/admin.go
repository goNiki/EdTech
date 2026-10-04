package auth

import (
	"context"
	"fmt"

	"edtech/internal/domain"
	errorsAPP "edtech/pkg/errors"

	"github.com/jackc/pgx/v5"
)

func (s *service) ChangeUserRole(ctx context.Context, adminID int64, targetUserID int64, newRole domain.Role) error {
	const op = "service.auth.ChangeUserRole"

	admin, err := s.repo.GetUserByID(ctx, adminID)
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

	if _, err := s.repo.GetUserByID(ctx, targetUserID); err != nil {
		return fmt.Errorf("%s: get target user: %w", op, err)
	}

	if err := s.repo.UpdateRole(ctx, targetUserID, newRole); err != nil {
		return fmt.Errorf("%s: update role: %w", op, err)
	}

	return nil
}

func (s *service) SetUserBanned(ctx context.Context, adminID int64, targetUserID int64, isBanned bool) error {
	const op = "service.auth.SetUserBanned"

	admin, err := s.repo.GetUserByID(ctx, adminID)
	if err != nil {
		return fmt.Errorf("%s: get admin: %w", op, err)
	}
	if !admin.IsAdmin() {
		return errorsAPP.ErrForbidden
	}

	if adminID == targetUserID {
		return errorsAPP.ErrCannotModifySelf
	}

	if _, err := s.repo.GetUserByID(ctx, targetUserID); err != nil {
		return fmt.Errorf("%s: get target user: %w", op, err)
	}

	err = s.txmanager.WithTX(ctx, pgx.TxOptions{}, func(ctx context.Context) error {
		if err := s.repo.SetBannedStatus(ctx, targetUserID, isBanned); err != nil {
			return err
		}

		if isBanned {
			if err := s.refreshRepo.DeleteAllByUserID(ctx, targetUserID); err != nil {
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

func (s *service) ListUsersForAdmin(ctx context.Context, adminID int64, filter domain.UserFilter, pagination domain.Pagination) ([]domain.User, int64, error) {
	const op = "service.auth.ListUsersForAdmin"

	admin, err := s.repo.GetUserByID(ctx, adminID)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: get admin: %w", op, err)
	}
	if !admin.IsAdmin() {
		return nil, 0, errorsAPP.ErrForbidden
	}

	users, total, err := s.repo.ListUsers(ctx, filter, pagination)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	return users, total, nil
}
