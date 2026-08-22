package auth

import (
	"fmt"
	"net/http"

	"edtech/internal/dto"
	"edtech/internal/infrastructure/logger"
	"edtech/internal/interfaces/handlers/converter"
	response "edtech/internal/interfaces/responce"
	errorsAPP "edtech/pkg/errors"

	"github.com/go-chi/render"
)

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.auth.Register"
	log := logger.GetLogger(r.Context(), op)

	var req dto.RegisterRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		response.HandleError(w, r, log, fmt.Errorf("%w: %v", errorsAPP.ErrDecodeJSON, err), op)
		return
	}

	if err := h.validator.Validate(req); err != nil {
		response.HandleError(w, r, log, fmt.Errorf("%w: %v", errorsAPP.ErrValidationFailed, err), op)
		return
	}

	user, err := h.authService.Register(r.Context(), converter.RegisterRequestToInput(req))
	if err != nil {
		response.HandleError(w, r, log, err, op)
		return
	}

	resp := dto.RegisterResponse{
		User: converter.UserToDTO(user),
	}

	response.Created(w, r, resp)
}
