package auth

import (
	"context"
	"edtech/internal/domain"
	"edtech/internal/dto"
	"edtech/internal/infrastructure/logger"
	"edtech/internal/infrastructure/logger/sl"
	errorsAPP "edtech/pkg/errors"

	"github.com/go-playground/validator/v10"
)

func (s *service) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	const op = "usercase.auth.Register"
	log := logger.GetLogger(ctx, op)

	validate := validator.New()

	if err := validate.Struct(req); err != nil {
		log.Error("failed validate fields", sl.Error(err))
		return nil, errorsAPP.ErrFailValidate
	}

	PassHash, err := s.hasherManager.Hash(req.Password)
	if err != nil {
		return nil, err
	}

	user := domain.User{
		Email:        req.Email,
		PasswordHash: PassHash,
		Username:     req.Username,
		Role:         domain.Role(req.Role),
	}

	newUser, err := s.repo.CreateUser(ctx, &user)
	if err != nil {
		return nil, err
	}

	resp := dto.RegisterResponse{
		ID:    newUser.ID,
		Email: newUser.Email,
	}

	return &resp, nil
}
