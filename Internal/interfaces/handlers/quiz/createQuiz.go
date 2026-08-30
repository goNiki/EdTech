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

// CreateQuiz handles POST /api/v1/lessons/{lesson_id}/quizzes
func (h *QuizHandler) CreateQuiz(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.quiz.CreateQuiz"
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

	var req dto.CreateQuizRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		response.HandleError(w, r, log, fmt.Errorf("%w: %w", errorsAPP.ErrDecodeJSON, err), op)
		return
	}

	if req.LessonID == 0 {
		req.LessonID = lessonID
	}

	if err := h.validator.Validate(req); err != nil {
		response.HandleError(w, r, log, fmt.Errorf("%w: %w", errorsAPP.ErrQuizValidation, err), op)
		return
	}

	quiz := domain.Quiz{
		LessonID:    req.LessonID,
		Title:       req.Title,
		Description: req.Description,
		PassingScor: req.PassingScor,
		MaxAttempts: req.MaxAttempts,
		TimeLimit:   req.TimeLimit,
	}

	createdQuiz, err := h.quizService.CreateQuiz(r.Context(), userID, &quiz)
	if err != nil {
		response.HandleError(w, r, log, err, op)
		return
	}

	response.Created(w, r, converter.QuizToDTO(createdQuiz))
}
