package auth

import (
	"context"
	"edtech/internal/domain"
	"edtech/internal/dto"
	"edtech/internal/infrastructure/hasher"
	"edtech/internal/infrastructure/jwt"
	"edtech/internal/infrastructure/logger"
	"edtech/internal/infrastructure/logger/sl"
	authRepo "edtech/internal/repository/auth"
	errorsAPP "edtech/pkg/errors"

	"github.com/go-playground/validator/v10"
)

type UserUseCase struct {
	repo       authRepo.UserRepository
	jwtManager jwt.TokenManager
	hasher     hasher.PasswordHasher
}

type UserService interface {
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error)
}

func NewUserUseCase(repo authRepo.UserRepository, jwtManager jwt.TokenManager, hasher hasher.PasswordHasher) *UserUseCase {
	return &UserUseCase{
		repo:       repo,
		jwtManager: jwtManager,
		hasher:     hasher,
	}
}

func (u *UserUseCase) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error) {

	log := logger.GetLogger(ctx)

	validate := validator.New()

	if err := validate.Struct(req); err != nil {
		log.Error("failed validate fields", sl.Error(err))
		return nil, errorsAPP.ErrFailValidate
	}

	PassHash, err := u.hasher.Hash(req.Password)
	if err != nil {
		return nil, err
	}

	var user domain.User

	user = domain.User{
		Email:        req.Email,
		PasswordHash: PassHash,
		Role:         domain.Role(req.Role),
	}

	newUser, err := u.repo.CreateUser(ctx, &user)
	if err != nil {
		return nil, err
	}

	resp := dto.RegisterResponse{
		ID:    newUser.ID,
		Email: newUser.Email,
	}

	return &resp, nil
}
