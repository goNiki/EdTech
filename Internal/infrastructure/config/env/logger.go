package env

import (
	errorsAPP "edtech/pkg/errors"
	"fmt"

	"github.com/caarlos0/env/v11"
)

type loggerEnvConfig struct {
	Level   string `env:"LOGGER_LEVEL,required"`
	Console consoleEnvConfig
	File    fileEnvConfig
}

type consoleEnvConfig struct {
	Enabled bool   `env:"LOGGER_CONSOLE_ENABLED,required"`
	Format  string `env:"LOGGER_CONSOLE_FORMAT,required"`
	Output  string `env:"LOGGER_CONSOLE_OUTPUT,required"`
}

type fileEnvConfig struct {
	Enabled  bool   `env:"LOGGER_FILE_ENABLED,required"`
	Format   string `env:"LOGGER_FILE_FORMAT,required"`
	FilePath string `env:"LOGGER_FILE_PATH,required"`
}

type ConsoleChannel interface {
	Enabled() bool
	Format() string
	Output() string
}

type FileChannel interface {
	Enabled() bool
	Format() string
	FilePath() string
}

type loggerConfig struct {
	raw loggerEnvConfig
}

type consoleConfig struct {
	raw consoleEnvConfig
}

type fileConfig struct {
	raw fileEnvConfig
}

func NewLoggerConfig() (*loggerConfig, error) {
	var raw loggerEnvConfig

	if err := env.Parse(&raw); err != nil {
		return nil, fmt.Errorf("%w: %w", errorsAPP.ErrParseLoggerConfig, err)
	}

	return &loggerConfig{
		raw: raw,
	}, nil
}

func (cfg *loggerConfig) Level() string {
	return cfg.raw.Level
}

func (cfg *loggerConfig) Console() ConsoleChannel {
	return &consoleConfig{
		raw: cfg.raw.Console,
	}
}

func (cfg *loggerConfig) File() FileChannel {
	return &fileConfig{
		raw: cfg.raw.File,
	}
}

func (cfg *consoleConfig) Enabled() bool {
	return cfg.raw.Enabled
}

func (cfg *consoleConfig) Format() string {
	return cfg.raw.Format
}

func (cfg *consoleConfig) Output() string {
	return cfg.raw.Output
}

func (cfg *fileConfig) Enabled() bool {
	return cfg.raw.Enabled
}

func (cfg *fileConfig) Format() string {
	return cfg.raw.Format
}

func (cfg *fileConfig) FilePath() string {
	return cfg.raw.FilePath
}
