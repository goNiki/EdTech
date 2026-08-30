package progress

import (
	"net/http"

	"edtech/internal/infrastructure/logger"
	"edtech/internal/interfaces/handlers/converter"
	response "edtech/internal/interfaces/response"
	errorsAPP "edtech/pkg/errors"
)

// GetLessonProgress handles GET /api/v1/lessons/{lesson_id}/progress
func (h *ProgressHandler) GetLessonProgress(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.progress.GetLessonProgress"
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

	prog, err := h.progressService.GetLessonProgress(r.Context(), userID, lessonID)
	if err != nil {
		response.HandleError(w, r, log, err, op)
		return
	}

	response.OK(w, r, converter.LessonProgressToDTO(prog))
}
