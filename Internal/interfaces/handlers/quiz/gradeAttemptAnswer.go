package quiz

import (
	"fmt"
	"net/http"

	"edtech/internal/dto"
	"edtech/internal/infrastructure/logger"
	"edtech/internal/interfaces/handlers/converter"
	response "edtech/internal/interfaces/response"
	errorsAPP "edtech/pkg/errors"

	"github.com/go-chi/render"
)

// GradeAttemptAnswer handles POST /api/v1/quizzes/attempts/{attempt_id}/answers/{answer_id}/grade
func (h *QuizHandler) GradeAttemptAnswer(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.quiz.GradeAttemptAnswer"
	log := logger.GetLogger(r.Context(), op)

	userID := h.getUserID(r.Context())
	if userID == 0 {
		response.HandleError(w, r, log, errorsAPP.ErrUnauthorized, op)
		return
	}

	attemptID, err := parseIDParam(r, "attempt_id", "attemptid", "id")
	if err != nil {
		response.HandleError(w, r, log, errorsAPP.ErrInvalidURLParam, op)
		return
	}

	answerID, err := parseIDParam(r, "answer_id", "answerid")
	if err != nil {
		response.HandleError(w, r, log, errorsAPP.ErrInvalidURLParam, op)
		return
	}

	var req dto.GradeAttemptRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		response.HandleError(w, r, log, fmt.Errorf("%w: %w", errorsAPP.ErrDecodeJSON, err), op)
		return
	}

	if err := h.validator.Validate(req); err != nil {
		response.HandleError(w, r, log, fmt.Errorf("%w: %w", errorsAPP.ErrValidationFailed, err), op)
		return
	}

	attempt, err := h.quizService.GradeAttemptAnswer(r.Context(), userID, attemptID, answerID, req.Points, req.Feedback)
	if err != nil {
		response.HandleError(w, r, log, err, op)
		return
	}

	response.OK(w, r, converter.QuizAttemptToDTO(attempt))
}
