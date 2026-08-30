package progress

import (
	"fmt"
	"net/http"

	"edtech/internal/domain"
	"edtech/internal/dto"
	"edtech/internal/infrastructure/logger"
	response "edtech/internal/interfaces/response"
	errorsAPP "edtech/pkg/errors"

	"github.com/go-chi/render"
)

// UpdateLessonProgress handles PATCH /api/v1/lessons/{lesson_id}/progress
func (h *ProgressHandler) UpdateLessonProgress(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.progress.UpdateLessonProgress"
	log := logger.GetLogger(r.Context(), op)

	userID := h.getUserID(r.Context())
	if userID == 0 {
		response.HandleError(w, r, log, errorsAPP.ErrUnauthorized, op)
		return
	}

	lessonID, err := parseIDParam(r, "lesson_id", "lessonId", "id")
	if err != nil {
		response.HandleError(w, r, log, errorsAPP.ErrInvalidURLParam, op)
		return
	}

	var req dto.UpdateLessonProgressRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		response.HandleError(w, r, log, fmt.Errorf("%w: %w", errorsAPP.ErrDecodeJSON, err), op)
		return
	}

	if err := h.validator.Validate(req); err != nil {
		response.HandleError(w, r, log, fmt.Errorf("%w: %w", errorsAPP.ErrValidationFailed, err), op)
		return
	}

	input := domain.UpdateProgressInput{
		Status:       domain.ProgressStatus(req.Status),
		TimeSpent:    req.TimeSpent,
		LastPosition: req.LastPos,
	}

	if err := h.progressService.UpdateLessonProgress(r.Context(), userID, lessonID, input); err != nil {
		response.HandleError(w, r, log, err, op)
		return
	}

	response.OK(w, r, map[string]string{"message": "progress updated"})
}
