package progress

import (
	"net/http"

	"edtech/internal/domain"
	"edtech/internal/dto"
	"edtech/internal/infrastructure/logger"
	response "edtech/internal/interfaces/response"
	errorsAPP "edtech/pkg/errors"

	"github.com/go-chi/render"
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

	var req dto.CompleteLessonRequest
	if r.ContentLength > 0 {
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Warn("failed to decode complete lesson request body", "error", err)
		}
	}

	essays := make([]domain.EssaySubmission, 0, len(req.Essays))
	for _, e := range req.Essays {
		essays = append(essays, domain.EssaySubmission{
			QuestionText: e.QuestionText,
			AnswerText:   e.AnswerText,
			MaxPoints:    e.MaxPoints,
		})
	}

	if err := h.progressService.CompleteLesson(r.Context(), userID, lessonID, req.Score, essays); err != nil {
		response.HandleError(w, r, log, err, op)
		return
	}

	response.OK(w, r, map[string]string{"message": "lesson completed"})
}
