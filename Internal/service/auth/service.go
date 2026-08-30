package auth

import (
	"edtech/internal/infrastructure/db"
	"edtech/internal/infrastructure/hasher"
	"edtech/internal/infrastructure/jwt"
	"edtech/internal/infrastructure/txmanager"
	"edtech/internal/repository"
)

type service struct {
	db            db.QueryExecutor
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
	database db.QueryExecutor,
	txmanager txmanager.TransactionManager,
) *service {
	return &service{
		db:            database,
		repo:          repo,
		jwtManager:    jwtManager,
		hasherManager: hasherManager,
		refreshRepo:   refreshRepo,
		txmanager:     txmanager,
	}
}
