package logger

import (
	"edtech/internal/infrastructure/logger"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

func (m *middlewareloger) Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		entry := m.logger.With(
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("remote_add", r.RemoteAddr),
			slog.String("user_agent", r.UserAgent()),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		ctx := logger.SetLogger(r.Context(), entry)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)

		entry.Info("Request finished")
	})
}
