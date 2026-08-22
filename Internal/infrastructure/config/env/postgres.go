package env

import (
	errorsAPP "edtech/pkg/errors"
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

type postgresEnvConfig struct {
	Host              string        `env:"DB_HOST,required"`
	Port              string        `env:"DB_PORT,required"`
	User              string        `env:"DB_USER,required"`
	Password          string        `env:"DB_PASSWORD,required"`
	Name              string        `env:"DB_NAME,required"`
	SslMode           string        `env:"DB_SSLMODE,required"`
	MaxConns          int32         `env:"DB_MAXCONNS,required"`
	MinConns          int32         `env:"DB_MINCONNS,required"`
	MaxConnLifeTime   time.Duration `env:"DB_MAXCONNLIFETIME,required"`
	MaxConnIdleTime   time.Duration `env:"DB_MAXCONNIDLETIME,required"`
	HealthCheckPeriod time.Duration `env:"DB_HEALTHCHECKPERIOD,required"`
}

type postgresConfig struct {
	raw postgresEnvConfig
}

func NewPostgresConfig() (*postgresConfig, error) {
	var raw postgresEnvConfig

	if err := env.Parse(&raw); err != nil {
		return nil, fmt.Errorf("%w: %w", errorsAPP.ErrParsePostgresConfig, err)
	}

	return &postgresConfig{
		raw: raw,
	}, nil

}

func Parse(raw interface{}) error {
	return nil
}

func (cfg *postgresConfig) Address() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", cfg.raw.User, cfg.raw.Password, cfg.raw.Host, cfg.raw.Port, cfg.raw.Name, cfg.raw.SslMode)
}

func (cfg *postgresConfig) Host() string {
	return cfg.raw.Host
}

func (cfg *postgresConfig) Port() string {
	return cfg.raw.Port
}

func (cfg *postgresConfig) User() string {
	return cfg.raw.User
}

func (cfg *postgresConfig) Password() string {
	return cfg.raw.Password
}

func (cfg *postgresConfig) Name() string {
	return cfg.raw.Name
}

func (cfg *postgresConfig) SslMode() string {
	return cfg.raw.SslMode
}

func (cfg *postgresConfig) MaxConns() int32 {
	return cfg.raw.MaxConns
}

func (cfg *postgresConfig) MinConns() int32 {
	return cfg.raw.MinConns
}

func (cfg *postgresConfig) MaxConnLifeTime() time.Duration {
	return cfg.raw.MaxConnLifeTime
}

func (cfg *postgresConfig) MaxConnIdleTime() time.Duration {
	return cfg.raw.MaxConnIdleTime
}

func (cfg *postgresConfig) HealthCheckPeriod() time.Duration {
	return cfg.raw.HealthCheckPeriod
}
