// Package config содержит функции и структуры для работы с конфигурацией приложения.
package config

import (
	"edtech/internal/infrastructure/config/env"
	errorsAPP "edtech/pkg/errors"
	"fmt"
	"os"

	"github.com/subosito/gotenv"
)

type config struct {
	Server   Server
	Postgres Postgres
	Logger   Logger
	JWT      JWT
}

func Load(path string) (*config, error) {

	if _, err := os.Stat(path); err == nil {
		if err := gotenv.Load(path); err != nil {
			return nil, fmt.Errorf("%w: %w", errorsAPP.ErrLoadEnv, err)
		}
	}

	server, err := env.NewServerConfig()
	if err != nil {
		return nil, err
	}

	postgres, err := env.NewPostgresConfig()
	if err != nil {
		return nil, err
	}

	logger, err := env.NewLoggerConfig()
	if err != nil {
		return nil, err
	}

	jwt, err := env.NewJwtConfig()
	if err != nil {
		return nil, err
	}

	return &config{
		Server:   server,
		Postgres: postgres,
		Logger:   logger,
		JWT:      jwt,
	}, nil

}
