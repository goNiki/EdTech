package progress

import (
	"net/http"

	"edtech/internal/dto"
	"edtech/internal/infrastructure/logger"
	response "edtech/internal/interfaces/response"
	errorsAPP "edtech/pkg/errors"

	"github.com/go-chi/render"
)

// SaveAttemptDraft handles PATCH /api/v1/lessons/{lesson_id}/attempts/{attempt_id}/draft
func (h *ProgressHandler) SaveAttemptDraft(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.progress.SaveAttemptDraft"
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

	attemptID, err := parseIDParam(r, "attempt_id", "attemptId")
	if err != nil {
		response.HandleError(w, r, log, errorsAPP.ErrInvalidURLParam, op)
		return
	}

	var req dto.SaveAttemptDraftRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		response.HandleError(w, r, log, errorsAPP.ErrDecodeJSON, op)
		return
	}

	savedAt, err := h.progressService.SaveAttemptDraft(r.Context(), userID, lessonID, attemptID, req.CurrentStep, req.Answers)
	if err != nil {
		response.HandleError(w, r, log, err, op)
		return
	}

	response.OK(w, r, dto.SaveAttemptDraftResponse{
		SavedAt: savedAt,
	})
}
