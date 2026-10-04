package auth

import (
	"fmt"
	"net/http"

	"edtech/internal/dto"
	"edtech/internal/infrastructure/logger"
	"edtech/internal/interfaces/handlers/converter"
	response "edtech/internal/interfaces/response"
	errorsAPP "edtech/pkg/errors"

	"github.com/go-chi/render"
)

func (h *AuthHandler) UpdatePreferences(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.auth.UpdatePreferences"
	log := logger.GetLogger(r.Context(), op)

	userID := h.authMiddleware.GetUserID(r.Context())
	if userID == 0 {
		response.HandleError(w, r, log, errorsAPP.ErrUnauthorized, op)
		return
	}

	var req dto.UpdatePreferencesRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		response.HandleError(w, r, log, fmt.Errorf("%w: %v", errorsAPP.ErrDecodeJSON, err), op)
		return
	}

	if err := h.validator.Validate(req); err != nil {
		response.HandleError(w, r, log, fmt.Errorf("%w: %v", errorsAPP.ErrValidationFailed, err), op)
		return
	}

	prefs, err := h.authService.UpdatePreferences(r.Context(), userID, converter.UpdatePreferencesRequestToDomain(req))
	if err != nil {
		response.HandleError(w, r, log, err, op)
		return
	}

	resp := dto.UpdatePreferencesResponse{
		Code:    http.StatusOK,
		Message: "preferences updated successfully",
		Data: dto.UpdatePreferencesData{
			Preferences: converter.PreferencesToDTO(prefs),
		},
	}

	response.OK(w, r, resp)
}
