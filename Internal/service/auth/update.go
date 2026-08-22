package auth

import (
	"context"
	"fmt"

	"edtech/internal/domain"
	errorsAPP "edtech/pkg/errors"
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

	if err := s.repo.UpdatePassword(ctx, s.db, user.ID, passHash); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
