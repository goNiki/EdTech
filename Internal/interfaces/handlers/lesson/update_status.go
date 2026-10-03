package lesson

import (
	"net/http"
	"strconv"

	"edtech/internal/dto"
	"edtech/internal/infrastructure/logger"
	response "edtech/internal/interfaces/response"
	errorsAPP "edtech/pkg/errors"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func (h *LessonHandler) UpdateLessonStatus(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.lesson.UpdateLessonStatus"

	log := logger.GetLogger(r.Context(), op)

	userID := h.authMiddleware.GetUserID(r.Context())
	if userID == 0 {
		response.HandleError(w, r, log, errorsAPP.ErrUnauthorized, op)
		return
	}

	lessonIDStr := chi.URLParam(r, "id")
	lessonID, err := strconv.ParseInt(lessonIDStr, 10, 64)
	if err != nil {
		response.HandleError(w, r, log, errorsAPP.ErrInvalidURLParam, op)
		return
	}

	var req dto.UpdateStatusRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		response.HandleError(w, r, log, errorsAPP.ErrDecodeJSON, op)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		response.HandleError(w, r, log, errorsAPP.ErrValidationFailed, op)
		return
	}

	if err := h.lessonService.UpdateLessonStatus(r.Context(), userID, lessonID, req.Status); err != nil {
		response.HandleError(w, r, log, err, op)
		return
	}

	response.OK(w, r, map[string]interface{}{
		"success": true,
		"message": "lesson status updated",
		"status":  req.Status,
	})
}
