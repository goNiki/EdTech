package auth

import (
	"context"
	"edtech/internal/infrastructure/jwt"
	"net/http"
)

type contextKey string

const (
	userIDKey   contextKey = "userID"
	userRoleKey contextKey = "role"
)

type middleware struct {
	jwtManager jwt.TokenManager
}

func NewAuthMiddleware(jwtManager jwt.TokenManager) *middleware {
	return &middleware{
		jwtManager: jwtManager,
	}
}

type AuthMiddleware interface {
	GetUserID(ctx context.Context) int64
	GetUserRole(ctx context.Context) string
	JWTMiddleware(next http.Handler) http.Handler
	RequireRole(allowedRoles []string) func(http.Handler) http.Handler
}

func (m *middleware) GetUserID(ctx context.Context) int64 {
	return GetUserID(ctx)
}

func (m *middleware) GetUserRole(ctx context.Context) string {
	return GetUserRole(ctx)
}

func GetUserID(ctx context.Context) int64 {
	userID, ok := ctx.Value(userIDKey).(int64)
	if !ok {
		return 0
	}
	return userID
}

func GetUserRole(ctx context.Context) string {
	userRole, ok := ctx.Value(userRoleKey).(string)
	if !ok {
		return ""
	}
	return userRole
}
