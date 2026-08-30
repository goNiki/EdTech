package progress

import (
	"net/http"

	"edtech/internal/infrastructure/logger"
	response "edtech/internal/interfaces/response"
	errorsAPP "edtech/pkg/errors"
)

// CompleteLesson handles POST /api/v1/lessons/{lesson_id}/complete
func (h *ProgressHandler) CompleteLesson(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.progress.CompleteLesson"
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

	if err := h.progressService.CompleteLesson(r.Context(), userID, lessonID); err != nil {
		response.HandleError(w, r, log, err, op)
		return
	}

	response.OK(w, r, map[string]string{"message": "lesson completed"})
}
