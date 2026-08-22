package auth

import (
	"edtech/internal/infrastructure/validator"
	"edtech/internal/interfaces/middleware/auth"
	"edtech/internal/service"
)

type AuthHandler struct {
	authService    service.AuthService
	authMiddleware auth.AuthMiddleware
	validator      validator.Validator
}

func NewAuthHandler(
	authService service.AuthService,
	authMiddleware auth.AuthMiddleware,
) *AuthHandler {
	return &AuthHandler{
		authService:    authService,
		authMiddleware: authMiddleware,
		validator:      *validator.NewValidator(),
	}
}
