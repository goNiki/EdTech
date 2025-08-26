package auth

import (
	"edtech/internal/dto"
	"edtech/internal/infrastructure/logger"
	"edtech/internal/infrastructure/logger/sl"
	authUC "edtech/internal/usecase/auth"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type AuthHandler struct {
	userService authUC.UserService
}

func NewAuthHandler(usecase authUC.UserService) *AuthHandler {
	return &AuthHandler{
		userService: usecase,
	}
}

func (a *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.auth.Register"

	log := logger.GetLogger(r.Context())

	logger := log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	logger.Info("register request started")

	var req dto.RegisterRequest

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		logger.Error("failed decode json", sl.Error(err))
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	resp, err := a.userService.Register(r.Context(), &req)
	if err != nil {
		logger.Error("failed service ", sl.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError) // сделать обработчик ошибок
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, dto.RegisterResponse{
		ID:    resp.ID,
		Email: resp.Email,
	})

}
