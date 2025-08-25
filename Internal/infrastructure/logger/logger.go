package logger

import (
	"context"
	"log/slog"
	"os"
)

type ctxKey string

const LoggerKey ctxKey = "logger"

func Init() *slog.Logger {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	log.Info("Logger initialized")

	return log
}

func GetLogger(ctx context.Context) *slog.Logger {
	l, ok := ctx.Value(LoggerKey).(*slog.Logger)
	if !ok {
		return slog.Default()

	}
	return l
}

func SetLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, LoggerKey, logger)
}
