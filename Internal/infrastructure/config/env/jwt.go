package env

import (
	errorsAPP "edtech/pkg/errors"
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

type jwtEnvConfig struct {
	Secret     string        `env:"JWT_SECRET,required"`
	AccessExp  time.Duration `env:"JWT_ACCESSEXP,required"`
	RefreshExp time.Duration `env:"JWT_REFRESHEXP,required"`
}

type jwtConfig struct {
	raw jwtEnvConfig
}

func NewJwtConfig() (*jwtConfig, error) {
	var raw jwtEnvConfig

	if err := env.Parse(&raw); err != nil {
		return nil, fmt.Errorf("%w: %w", errorsAPP.ErrParseJWTConfig, err)
	}

	return &jwtConfig{
		raw: raw,
	}, nil
}

func (cfg *jwtConfig) Secret() string {
	return cfg.raw.Secret
}

func (cfg *jwtConfig) AccessExp() time.Duration {
	return cfg.raw.AccessExp
}

func (cfg *jwtConfig) RefreshExp() time.Duration {
	return cfg.raw.RefreshExp
}
