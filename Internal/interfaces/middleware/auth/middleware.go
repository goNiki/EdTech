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

type AuthMiddleware interface {
	GetUserID(ctx context.Context) int64
	GetUserRole(ctx context.Context) string
}

func (m *middleware) GetUserID(ctx context.Context) int64 {
	return ctx.Value(userIDKey).(int64)
}

func (m *middleware) GetUserRole(ctx context.Context) string {
	return ctx.Value(userRoleKey).(string)
}
