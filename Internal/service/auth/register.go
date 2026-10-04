package auth

import (
	"context"
	"edtech/internal/domain"
	"fmt"
)

func (s *service) Register(ctx context.Context, input domain.RegisterInput) (*domain.User, error) {
	const op = "usercase.auth.Register"

	passHash, err := s.hasherManager.Hash(input.Password)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	createUser := domain.CreateUser{
		Email:        input.Email,
		PasswordHash: passHash,
		Username:     input.Username,
		Role:         input.Role,
	}

	newUser, err := s.repo.CreateUser(ctx, createUser)
	if err != nil {
		// we already handled ErrEmailAlreadyExists and ErrUserNameAlreadyExists in the repo
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return newUser, nil
}
