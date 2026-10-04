package env

import (
	errorsAPP "edtech/pkg/errors"
	"fmt"
	"net"
	"time"

	"github.com/caarlos0/env/v11"
)

type serverEnvConfig struct {
	Host               string        `env:"SERVER_HOST,required"`
	Port               string        `env:"SERVER_PORT,required"`
	TimeOut            time.Duration `env:"SERVER_TIMEOUT,required"`
	Idletimeout        time.Duration `env:"SERVER_IDLETIMEOUT,required"`
	ReadHeaderTimeout  time.Duration `env:"SERVER_READ_HEADER_TIMEOUT" envDefault:"5s"`
	CorsAllowedOrigins []string      `env:"CORS_ALLOWED_ORIGINS" envSeparator:"," envDefault:"http://localhost:3000,http://localhost:3001,http://127.0.0.1:3000"`
}

type serverConfig struct {
	raw serverEnvConfig
}

func NewServerConfig() (*serverConfig, error) {
	var raw serverEnvConfig

	if err := env.Parse(&raw); err != nil {
		return nil, fmt.Errorf("%w: %w", errorsAPP.ErrParseServerConfig, err)
	}

	return &serverConfig{
		raw: raw,
	}, nil
}

func (cfg *serverConfig) Host() string {
	return cfg.raw.Host
}

func (cfg *serverConfig) Port() string {
	return cfg.raw.Port
}

func (cfg *serverConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}

func (cfg *serverConfig) TimeOut() time.Duration {
	return cfg.raw.TimeOut
}

func (cfg *serverConfig) Idletimeout() time.Duration {
	return cfg.raw.Idletimeout
}

func (cfg *serverConfig) ReadHeaderTimeout() time.Duration {
	return cfg.raw.ReadHeaderTimeout
}

func (cfg *serverConfig) CorsAllowedOrigins() []string {
	return cfg.raw.CorsAllowedOrigins
}
