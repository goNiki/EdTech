package notification

import (
	"log/slog"
	"net/http"
	"strconv"

	"edtech/internal/dto"
	"edtech/internal/interfaces/middleware/auth"
	"edtech/internal/interfaces/response"
	"edtech/internal/service"
	errorsAPP "edtech/pkg/errors"

	"github.com/go-chi/chi/v5"
)

type NotificationHandler struct {
	notificationService service.NotificationServices
	log                 *slog.Logger
	authMiddleware      auth.AuthMiddleware
}

func NewNotificationHandler(
	notificationService service.NotificationServices,
	log *slog.Logger,
	authMiddleware auth.AuthMiddleware,
) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
		log:                 log,
		authMiddleware:      authMiddleware,
	}
}

// GetNotifications handles GET /api/v1/notifications
func (h *NotificationHandler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.notification.GetNotifications"

	userID := h.authMiddleware.GetUserID(r.Context())
	if userID == 0 {
		response.HandleError(w, r, h.log, errorsAPP.ErrUnauthorized, op)
		return
	}

	limit := 20
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			if parsedLimit > 100 {
				parsedLimit = 100
			}
			limit = parsedLimit
		}
	}

	feed, err := h.notificationService.GetFeed(r.Context(), userID, limit)
	if err != nil {
		response.HandleError(w, r, h.log, err, op)
		return
	}

	response.OK(w, r, dto.ToNotificationFeedResponse(feed))
}

// MarkAsRead handles PATCH /api/v1/notifications/{id}/read
func (h *NotificationHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.notification.MarkAsRead"

	userID := h.authMiddleware.GetUserID(r.Context())
	if userID == 0 {
		response.HandleError(w, r, h.log, errorsAPP.ErrUnauthorized, op)
		return
	}

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.HandleError(w, r, h.log, errorsAPP.ErrInvalidURLParam, op)
		return
	}

	if err := h.notificationService.MarkAsRead(r.Context(), userID, id); err != nil {
		response.HandleError(w, r, h.log, err, op)
		return
	}

	response.OK(w, r, map[string]string{"message": "notification marked as read"})
}

// MarkAllAsRead handles POST /api/v1/notifications/read-all
func (h *NotificationHandler) MarkAllAsRead(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.notification.MarkAllAsRead"

	userID := h.authMiddleware.GetUserID(r.Context())
	if userID == 0 {
		response.HandleError(w, r, h.log, errorsAPP.ErrUnauthorized, op)
		return
	}

	if err := h.notificationService.MarkAllAsRead(r.Context(), userID); err != nil {
		response.HandleError(w, r, h.log, err, op)
		return
	}

	response.OK(w, r, map[string]string{"message": "all notifications marked as read"})
}
