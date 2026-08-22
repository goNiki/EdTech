package auth

import (
	"fmt"
	"net/http"

	"edtech/internal/dto"
	"edtech/internal/infrastructure/logger"
	response "edtech/internal/interfaces/responce"
	errorsAPP "edtech/pkg/errors"

	"github.com/go-chi/render"
)

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.auth.Logout"
	log := logger.GetLogger(r.Context(), op)

	var req dto.LogoutRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		response.HandleError(w, r, log, fmt.Errorf("%w: %v", errorsAPP.ErrDecodeJSON, err), op)
		return
	}

	if err := h.validator.Validate(req); err != nil {
		response.HandleError(w, r, log, fmt.Errorf("%w: %v", errorsAPP.ErrValidationFailed, err), op)
		return
	}

	if err := h.authService.Logout(r.Context(), req.RefreshToken); err != nil {
		response.HandleError(w, r, log, err, op)
		return
	}

	response.OK(w, r, map[string]string{"message": "successfully logged out"})
}
