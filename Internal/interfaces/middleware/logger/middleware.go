package logger

import "log/slog"

type middleware struct {
	logger *slog.Logger
}

func NewLoggerMiddleware(logger *slog.Logger) *middleware {
	return &middleware{
		logger: logger,
	}
}
