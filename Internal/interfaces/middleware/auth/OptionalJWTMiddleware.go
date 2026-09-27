package auth

import (
	"context"
	"net/http"
	"strings"
)

func (m *middleware) OptionalJWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			next.ServeHTTP(w, r)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		id, role, err := m.jwtManager.ParseToken(tokenStr)
		if err == nil {
			ctx := context.WithValue(r.Context(), userIDKey, id)
			ctx = context.WithValue(ctx, userRoleKey, role)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}
