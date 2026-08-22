package auth

import (
	"context"
	"edtech/internal/domain"
	"fmt"
)

func (s *service) GetCurrentUser(ctx context.Context, userID int64) (*domain.User, error) {
	const op = "service.auth.GetCurrentUser"

	user, err := s.repo.GetUserByID(ctx, s.db, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil

}
