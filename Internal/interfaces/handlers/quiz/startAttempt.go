package quiz

import (
	"net/http"

	"edtech/internal/infrastructure/logger"
	"edtech/internal/interfaces/handlers/converter"
	response "edtech/internal/interfaces/response"
	errorsAPP "edtech/pkg/errors"
)

// StartAttempt handles POST /api/v1/quizzes/{quiz_id}/attempts/start
func (h *QuizHandler) StartAttempt(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.quiz.StartAttempt"
	log := logger.GetLogger(r.Context(), op)

	userID := h.getUserID(r.Context())
	if userID == 0 {
		response.HandleError(w, r, log, errorsAPP.ErrUnauthorized, op)
		return
	}

	quizID, err := parseIDParam(r, "quiz_id", "quizid", "id")
	if err != nil {
		response.HandleError(w, r, log, errorsAPP.ErrInvalidURLParam, op)
		return
	}

	attempt, err := h.quizService.StartAttempt(r.Context(), userID, quizID)
	if err != nil {
		response.HandleError(w, r, log, err, op)
		return
	}

	response.Created(w, r, converter.QuizAttemptToDTO(attempt))
}
