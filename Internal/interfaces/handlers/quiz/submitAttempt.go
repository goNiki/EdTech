package quiz

import (
	"fmt"
	"net/http"

	"edtech/internal/domain"
	"edtech/internal/dto"
	"edtech/internal/infrastructure/logger"
	"edtech/internal/interfaces/handlers/converter"
	response "edtech/internal/interfaces/response"
	errorsAPP "edtech/pkg/errors"

	"github.com/go-chi/render"
)

// SubmitAttempt handles POST /api/v1/quizzes/attempts/{attempt_id}/submit
func (h *QuizHandler) SubmitAttempt(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.quiz.SubmitAttempt"
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

	var req dto.SubmitAttemptRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		response.HandleError(w, r, log, fmt.Errorf("%w: %w", errorsAPP.ErrDecodeJSON, err), op)
		return
	}

	if err := h.validator.Validate(req); err != nil {
		response.HandleError(w, r, log, fmt.Errorf("%w: %w", errorsAPP.ErrValidationFailed, err), op)
		return
	}

	answers := make([]domain.QuizAttemptAnswer, 0, len(req.Answers))
	for _, a := range req.Answers {
		textVal := ""
		if a.TextValue != nil {
			textVal = *a.TextValue
		}
		answers = append(answers, domain.QuizAttemptAnswer{
			AttemptID:  attemptID,
			QuestionID: a.QuestionID,
			AnswerID:   a.AnswerID,
			TextValue:  textVal,
		})
	}

	attempt, err := h.quizService.SubmitAttempt(r.Context(), userID, attemptID, answers)
	if err != nil {
		response.HandleError(w, r, log, err, op)
		return
	}

	response.OK(w, r, converter.QuizAttemptToDTO(attempt))
}
