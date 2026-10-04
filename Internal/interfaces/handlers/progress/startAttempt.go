package progress

import (
	"net/http"

	"edtech/internal/dto"
	"edtech/internal/infrastructure/logger"
	response "edtech/internal/interfaces/response"
	errorsAPP "edtech/pkg/errors"
)

// StartAttempt handles POST /api/v1/lessons/{lesson_id}/attempts/start
func (h *ProgressHandler) StartAttempt(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.progress.StartAttempt"
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

	result, err := h.progressService.StartLessonAttempt(r.Context(), userID, lessonID)
	if err != nil {
		response.HandleError(w, r, log, err, op)
		return
	}

	response.OK(w, r, dto.StartAttemptResponse{
		AttemptID: result.AttemptID,
		StartedAt: result.StartedAt,
	})
}
