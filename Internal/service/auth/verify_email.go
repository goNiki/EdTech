package auth

import (
	"context"
	"fmt"
)

func (s *service) VerifyEmail(ctx context.Context, userID int64) error {
	const op = "service.auth.VerifyEmail"

	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("%s: get user: %w", op, err)
	}

	if user.EmailVerified {
		return nil
	}

	if err := s.repo.SetEmailVerified(ctx, userID, true); err != nil {
		return fmt.Errorf("%s: set email verified: %w", op, err)
	}

	return nil
}
