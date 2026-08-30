package auth

import (
	"net/http"

	"edtech/internal/infrastructure/logger"
	response "edtech/internal/interfaces/response"
	errorsAPP "edtech/pkg/errors"
)

func (h *AuthHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.auth.VerifyEmail"
	log := logger.GetLogger(r.Context(), op)

	userID := h.authMiddleware.GetUserID(r.Context())
	if userID == 0 {
		response.HandleError(w, r, log, errorsAPP.ErrUnauthorized, op)
		return
	}

	if err := h.authService.VerifyEmail(r.Context(), userID); err != nil {
		response.HandleError(w, r, log, err, op)
		return
	}

	response.OK(w, r, map[string]string{"message": "email successfully verified"})
}
