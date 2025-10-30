package auth

import (
	"edtech/internal/infrastructure/hasher"
	"edtech/internal/infrastructure/jwt"
	"edtech/internal/repository"
)

type service struct {
	repo          repository.UserRepository
	jwtManager    jwt.TokenManager
	hasherManager hasher.HasherManager
	refreshRepo   repository.RefreshRepository
}

func NewAuthService(repo repository.UserRepository, jwtManager jwt.TokenManager, hasherManager hasher.HasherManager, refreshRepo repository.RefreshRepository) *service {
	return &service{
		repo:          repo,
		jwtManager:    jwtManager,
		hasherManager: hasherManager,
		refreshRepo:   refreshRepo,
	}
}
