package progress

import (
	"net/http"

	"edtech/internal/dto"
	"edtech/internal/infrastructure/logger"
	response "edtech/internal/interfaces/response"
	errorsAPP "edtech/pkg/errors"
)

// GetActiveAttempt handles GET /api/v1/lessons/{lesson_id}/attempts/active
func (h *ProgressHandler) GetActiveAttempt(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.progress.GetActiveAttempt"
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

	activeResult, err := h.progressService.GetActiveLessonAttempt(r.Context(), userID, lessonID)
	if err != nil {
		response.HandleError(w, r, log, err, op)
		return
	}

	if !activeResult.HasActiveAttempt || activeResult.Attempt == nil {
		response.OK(w, r, dto.ActiveAttemptResponse{
			HasActiveAttempt: false,
		})
		return
	}

	response.OK(w, r, dto.ActiveAttemptResponse{
		HasActiveAttempt: true,
		Attempt: &dto.ActiveAttemptDTO{
			ID:               activeResult.Attempt.ID,
			StartedAt:        activeResult.Attempt.StartedAt,
			TimeLimitMinutes: activeResult.Attempt.TimeLimitMinutes,
			RemainingSeconds: activeResult.Attempt.RemainingSeconds,
			CurrentStep:      activeResult.Attempt.CurrentStep,
			DraftAnswers:     activeResult.Attempt.DraftAnswers,
		},
	})
}
