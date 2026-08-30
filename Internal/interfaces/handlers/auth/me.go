package auth

import (
	"net/http"

	"edtech/internal/infrastructure/logger"
	"edtech/internal/interfaces/handlers/converter"
	response "edtech/internal/interfaces/response"
	errorsAPP "edtech/pkg/errors"
)

func (h *AuthHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.auth.GetCurrentUser"
	log := logger.GetLogger(r.Context(), op)

	userID := h.authMiddleware.GetUserID(r.Context())
	if userID == 0 {
		response.HandleError(w, r, log, errorsAPP.ErrUnauthorized, op)
		return
	}

	user, err := h.authService.GetCurrentUser(r.Context(), userID)
	if err != nil {
		response.HandleError(w, r, log, err, op)
		return
	}

	response.OK(w, r, converter.UserToDTO(user))
}
