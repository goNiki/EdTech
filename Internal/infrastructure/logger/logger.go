package logger

import (
	"context"
	"log/slog"
	"os"

	"github.com/go-chi/chi/v5/middleware"
)

type ctxKey string

const LoggerKey ctxKey = "logger"

func Init() *slog.Logger {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	log.Info("Logger initialized")

	return log
}

func GetLogger(ctx context.Context, op string) *slog.Logger {
	l, ok := ctx.Value(LoggerKey).(*slog.Logger)
	if !ok || l == nil {
		l = slog.Default()
	}

	// создаём новый логгер с дополнительными полями
	return l.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(ctx)),
	)
}

func SetLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, LoggerKey, logger)
}
