package config

import (
	"edtech/internal/infrastructure/config/env"
	"time"
)

type Logger interface {
	Level() string
	Console() env.ConsoleChannel
	File() env.FileChannel
}

type Postgres interface {
	Address() string
	Host() string
	Port() string
	User() string
	Password() string
	Name() string
	SslMode() string
	MaxConns() int32
	MinConns() int32
	MaxConnLifeTime() time.Duration
	MaxConnIdleTime() time.Duration
	HealthCheckPeriod() time.Duration
}

type Server interface {
	Host() string
	Port() string
	Address() string
	TimeOut() time.Duration
	Idletimeout() time.Duration
}

type JWT interface {
	Secret() string
	AccessExp() time.Duration
	RefreshExp() time.Duration
}
