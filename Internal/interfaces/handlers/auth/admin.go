package auth

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"edtech/internal/domain"
	"edtech/internal/dto"
	"edtech/internal/infrastructure/logger"
	response "edtech/internal/interfaces/response"
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

func (h *AuthHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.auth.ListUsers"
	log := logger.GetLogger(r.Context(), op)

	adminID := h.authMiddleware.GetUserID(r.Context())
	if adminID == 0 {
		response.HandleError(w, r, log, errorsAPP.ErrUnauthorized, op)
		return
	}

	query := r.URL.Query()

	var filter domain.UserFilter
	if search := query.Get("search"); search != "" {
		filter.Search = &search
	}
	if roleStr := query.Get("role"); roleStr != "" {
		role := domain.Role(roleStr)
		filter.Role = &role
	}
	if isBannedStr := query.Get("is_banned"); isBannedStr != "" {
		if isBanned, err := strconv.ParseBool(isBannedStr); err == nil {
			filter.IsBanned = &isBanned
		}
	}

	page := int64(1)
	if pageStr := query.Get("page"); pageStr != "" {
		if p, err := strconv.ParseInt(pageStr, 10, 64); err == nil && p > 0 {
			page = p
		}
	}

	pageSize := int64(20)
	if pageSizeStr := query.Get("page_size"); pageSizeStr != "" {
		if ps, err := strconv.ParseInt(pageSizeStr, 10, 64); err == nil && ps > 0 {
			pageSize = ps
		}
	}

	pagination := domain.Pagination{
		Page:     page,
		PageSize: pageSize,
	}

	users, total, err := h.authService.ListUsersForAdmin(r.Context(), adminID, filter, pagination)
	if err != nil {
		response.HandleError(w, r, log, err, op)
		return
	}

	items := make([]dto.AdminUserItemResponse, 0, len(users))
	for _, u := range users {
		items = append(items, dto.AdminUserItemResponse{
			ID:        u.ID,
			Email:     u.Email,
			Username:  u.Username,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Role:      u.Role,
			IsBanned:  u.IsBanned,
			CreatedAt: u.CreatedAt.UTC().Format(time.RFC3339),
		})
	}

	pagination.Sanitize()
	response.OK(w, r, dto.AdminUsersListResponse{
		Users:    items,
		Total:    total,
		Page:     pagination.Page,
		PageSize: pagination.PageSize,
	})
}

