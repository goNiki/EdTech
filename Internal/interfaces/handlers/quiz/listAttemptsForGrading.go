package quiz

import (
	"net/http"
	"strconv"

	"edtech/internal/dto"
	"edtech/internal/infrastructure/logger"
	"edtech/internal/interfaces/handlers/converter"
	response "edtech/internal/interfaces/response"
	errorsAPP "edtech/pkg/errors"
)

// ListAttemptsForGrading handles GET /api/v1/courses/{course_id}/quizzes/attempts
func (h *QuizHandler) ListAttemptsForGrading(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.quiz.ListAttemptsForGrading"
	log := logger.GetLogger(r.Context(), op)

	userID := h.getUserID(r.Context())
	if userID == 0 {
		response.HandleError(w, r, log, errorsAPP.ErrUnauthorized, op)
		return
	}

	courseID, err := parseIDParam(r, "course_id", "courseid", "id")
	if err != nil {
		response.HandleError(w, r, log, errorsAPP.ErrInvalidURLParam, op)
		return
	}

	page := 1
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	pageSize := 10
	if sizeStr := r.URL.Query().Get("page_size"); sizeStr != "" {
		if s, err := strconv.Atoi(sizeStr); err == nil && s > 0 {
			pageSize = s
		}
	} else if sizeStr := r.URL.Query().Get("pageSize"); sizeStr != "" {
		if s, err := strconv.Atoi(sizeStr); err == nil && s > 0 {
			pageSize = s
		}
	} else if sizeStr := r.URL.Query().Get("limit"); sizeStr != "" {
		if s, err := strconv.Atoi(sizeStr); err == nil && s > 0 {
			pageSize = s
		}
	}

	var quizIDPtr *int64
	if quizIDStr := r.URL.Query().Get("quiz_id"); quizIDStr != "" {
		if qID, err := strconv.ParseInt(quizIDStr, 10, 64); err == nil && qID > 0 {
			quizIDPtr = &qID
		}
	} else if quizIDStr := r.URL.Query().Get("quizId"); quizIDStr != "" {
		if qID, err := strconv.ParseInt(quizIDStr, 10, 64); err == nil && qID > 0 {
			quizIDPtr = &qID
		}
	}

	attempts, total, err := h.quizService.ListAttemptsForGrading(r.Context(), userID, courseID, quizIDPtr, page, pageSize)
	if err != nil {
		response.HandleError(w, r, log, err, op)
		return
	}

	resp := dto.ListAttemptsResponse{
		Attempts: converter.QuizAttemptsToDTO(attempts),
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}

	response.OK(w, r, resp)
}
