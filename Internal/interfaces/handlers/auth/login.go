package auth

import (
	"fmt"
	"net/http"

	"edtech/internal/dto"
	"edtech/internal/infrastructure/logger"
	response "edtech/internal/interfaces/response"
	errorsAPP "edtech/pkg/errors"

	"github.com/go-chi/render"
)

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.auth.Login"
	log := logger.GetLogger(r.Context(), op)

	var req dto.LoginRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		response.HandleError(w, r, log, fmt.Errorf("%w: %v", errorsAPP.ErrDecodeJSON, err), op)
		return
	}

	if err := h.validator.Validate(req); err != nil {
		response.HandleError(w, r, log, fmt.Errorf("%w: %v", errorsAPP.ErrValidationFailed, err), op)
		return
	}

	authTokens, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		response.HandleError(w, r, log, err, op)
		return
	}

	resp := dto.LoginResponse{
		ID:           authTokens.ID,
		AccessToken:  authTokens.AccessToken,
		RefreshToken: authTokens.RefreshToken,
		ExpiresIn:    authTokens.ExpiresIn,
	}

	response.OK(w, r, resp)
}
