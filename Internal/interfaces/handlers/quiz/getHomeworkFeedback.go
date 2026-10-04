package quiz

import (
	"net/http"

	"edtech/internal/infrastructure/logger"
	"edtech/internal/interfaces/handlers/converter"
	response "edtech/internal/interfaces/response"
	errorsAPP "edtech/pkg/errors"
)

// GetStudentHomeworkFeedback handles GET /api/v1/lessons/{id}/homework-feedback
func (h *QuizHandler) GetStudentHomeworkFeedback(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.quiz.GetStudentHomeworkFeedback"
	log := logger.GetLogger(r.Context(), op)

	userID := h.getUserID(r.Context())
	if userID == 0 {
		response.HandleError(w, r, log, errorsAPP.ErrUnauthorized, op)
		return
	}

	lessonID, err := parseIDParam(r, "id", "lesson_id")
	if err != nil {
		response.HandleError(w, r, log, errorsAPP.ErrInvalidURLParam, op)
		return
	}

	feedback, err := h.quizService.GetStudentHomeworkFeedback(r.Context(), userID, lessonID)
	if err != nil {
		response.HandleError(w, r, log, err, op)
		return
	}

	response.OK(w, r, converter.StudentHomeworkFeedbackToDTO(feedback))
}
