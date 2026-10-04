package auth

import (
	"edtech/internal/infrastructure/hasher"
	"edtech/internal/infrastructure/jwt"
	"edtech/internal/infrastructure/txmanager"
	"edtech/internal/repository"
)

type service struct {
	txmanager     txmanager.TransactionManager
	repo          repository.UserRepository
	jwtManager    jwt.TokenManager
	hasherManager hasher.HasherManager
	refreshRepo   repository.RefreshRepository
}

func NewAuthService(
	repo repository.UserRepository,
	jwtManager jwt.TokenManager,
	hasherManager hasher.HasherManager,
	refreshRepo repository.RefreshRepository,
	txmanager txmanager.TransactionManager,
) *service {
	return &service{
		repo:          repo,
		jwtManager:    jwtManager,
		hasherManager: hasherManager,
		refreshRepo:   refreshRepo,
		txmanager:     txmanager,
	}
}
