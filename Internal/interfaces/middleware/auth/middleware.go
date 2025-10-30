package auth

import (
	"context"
	"edtech/internal/infrastructure/jwt"
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

func GetUserID(ctx context.Context) int64 {
	if v, ok := ctx.Value(userIDKey).(int64); ok {
		return v
	}
	return 0
}

func GetUserRole(ctx context.Context) string {
	if v, ok := ctx.Value(userRoleKey).(string); ok {
		return v
	}
	return ""
}
