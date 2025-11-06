package logger

import "log/slog"

type middlewareloger struct {
	logger *slog.Logger
}

func NewLoggerMiddleware(logger *slog.Logger) *middlewareloger {
	return &middlewareloger{
		logger: logger,
	}
}
