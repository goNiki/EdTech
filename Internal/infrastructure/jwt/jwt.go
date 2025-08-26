package jwt

import (
	"edtech/internal/infrastructure/config"
	"time"
)

type JwtManager struct {
	Secret     string
	AccessExp  time.Duration
	RefreshExp time.Duration
}

type TokenManager interface {
}

func NewJwtManager(cfg *config.JWTConfig) *JwtManager {
	return &JwtManager{
		Secret:     cfg.Secret,
		AccessExp:  cfg.AccesExp,
		RefreshExp: cfg.RefreshExp,
	}
}
