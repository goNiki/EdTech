package auth

import (
	"fmt"
	"net/http"
	"strconv"

	"edtech/internal/domain"
	"edtech/internal/dto"
	"edtech/internal/infrastructure/logger"
	response "edtech/internal/interfaces/responce"
	errorsAPP "edtech/pkg/errors"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func (h *AuthHandler) ChangeUserRole(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.auth.ChangeUserRole"
	log := logger.GetLogger(r.Context(), op)

	adminID := h.authMiddleware.GetUserID(r.Context())
	if adminID == 0 {
		response.HandleError(w, r, log, errorsAPP.ErrUnauthorized, op)
		return
	}

	targetIDStr := chi.URLParam(r, "id")
	targetID, err := strconv.ParseInt(targetIDStr, 10, 64)
	if err != nil {
		response.HandleError(w, r, log, errorsAPP.ErrInvalidURLParam, op)
		return
	}

	var req dto.ChangeUserRoleRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		response.HandleError(w, r, log, fmt.Errorf("%w: %v", errorsAPP.ErrDecodeJSON, err), op)
		return
	}

	if err := h.validator.Validate(req); err != nil {
		response.HandleError(w, r, log, fmt.Errorf("%w: %v", errorsAPP.ErrValidationFailed, err), op)
		return
	}

	if err := h.authService.ChangeUserRole(r.Context(), adminID, targetID, domain.Role(req.Role)); err != nil {
		response.HandleError(w, r, log, err, op)
		return
	}

	response.OK(w, r, map[string]string{"message": "user role successfully updated"})
}

func (h *AuthHandler) SetUserBanned(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.auth.SetUserBanned"
	log := logger.GetLogger(r.Context(), op)

	adminID := h.authMiddleware.GetUserID(r.Context())
	if adminID == 0 {
		response.HandleError(w, r, log, errorsAPP.ErrUnauthorized, op)
		return
	}

	targetIDStr := chi.URLParam(r, "id")
	targetID, err := strconv.ParseInt(targetIDStr, 10, 64)
	if err != nil {
		response.HandleError(w, r, log, errorsAPP.ErrInvalidURLParam, op)
		return
	}

	var req dto.SetUserBannedRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		response.HandleError(w, r, log, fmt.Errorf("%w: %v", errorsAPP.ErrDecodeJSON, err), op)
		return
	}

	if err := h.authService.SetUserBanned(r.Context(), adminID, targetID, req.IsBanned); err != nil {
		response.HandleError(w, r, log, err, op)
		return
	}

	msg := "user successfully unbanned"
	if req.IsBanned {
		msg = "user successfully banned"
	}

	response.OK(w, r, map[string]string{"message": msg})
}
