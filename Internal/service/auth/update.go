package auth

import (
	"context"
	"fmt"

	"edtech/internal/domain"
	"edtech/internal/infrastructure/db"
	errorsAPP "edtech/pkg/errors"

	"github.com/jackc/pgx/v5"
)

func (s *service) UpdateProfile(ctx context.Context, userID int64, input domain.UpdateProfileInput) (*domain.User, error) {
	const op = "service.auth.UpdateProfile"

	user, err := s.repo.GetUserByID(ctx, s.db, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: get user: %w", op, err)
	}

	user.UpdateProfile(input)

	if err := s.repo.UpdateProfile(ctx, s.db, user); err != nil {
		return nil, fmt.Errorf("%s: update profile: %w", op, err)
	}

	return user, nil
}

func (s *service) ChangePassword(ctx context.Context, userID int64, oldPassword, newPassword string) error {
	const op = "service.auth.ChangePassword"

	if len(newPassword) < 8 {
		return fmt.Errorf("%s: %w", op, errorsAPP.ErrPasswordTooShort)
	}
	if oldPassword == newPassword {
		return fmt.Errorf("%s: %w", op, errorsAPP.ErrSamePassword)
	}

	user, err := s.repo.GetUserByID(ctx, s.db, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	ok, err := s.hasherManager.CheckPassword(user.PasswordHash, oldPassword)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if !ok {
		return fmt.Errorf("%s: %w", op, errorsAPP.ErrInvalidCredentials)
	}

	passHash, err := s.hasherManager.Hash(newPassword)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = s.txmanager.WithTX(ctx, pgx.TxOptions{}, func(ctx context.Context, tx db.QueryExecutor) error {
		if err := s.repo.UpdatePassword(ctx, tx, user.ID, passHash); err != nil {
			return fmt.Errorf("%s: update password: %w", op, err)
		}
		if err := s.refreshRepo.DeleteAllByUserID(ctx, tx, user.ID); err != nil {
			return fmt.Errorf("%s: revoke refresh tokens: %w", op, err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
