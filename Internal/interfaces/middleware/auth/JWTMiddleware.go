package auth

import (
	"context"
	"net/http"
	"strings"
)

func (m *middleware) JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		// TODO обработчик ошибок
		id, role, err := m.jwtManager.ParseToken(tokenStr)
		if err != nil {
			http.Error(w, "token invalid", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, id)
		ctx = context.WithValue(ctx, userRoleKey, role)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
