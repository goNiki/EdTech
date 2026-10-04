package auth

import (
	"context"
	errorsAPP "edtech/pkg/errors"
	"errors"
	"fmt"
)

func (s *service) Logout(ctx context.Context, refreshToken string) error {
	const op = "service.auth.Logout"

	hashRefreshToken := s.hasherManager.HashRefreshToken(refreshToken)

	err := s.refreshRepo.DeleteRefreshToken(ctx, hashRefreshToken)
	if err != nil && !errors.Is(err, errorsAPP.ErrRefreshTokenNotFound) {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
