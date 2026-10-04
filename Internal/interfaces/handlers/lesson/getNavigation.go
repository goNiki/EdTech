package lesson

import (
	"net/http"
	"strconv"

	"edtech/internal/interfaces/handlers/converter"
	"edtech/internal/interfaces/response"
	errorsAPP "edtech/pkg/errors"

	"github.com/go-chi/chi/v5"
)

// GetLessonNavigationContext handles GET /api/v1/lessons/{id}/navigation
func (h *LessonHandler) GetLessonNavigationContext(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.lesson.GetLessonNavigationContext"

	idStr := chi.URLParam(r, "id")
	lessonID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.HandleError(w, r, h.log, errorsAPP.ErrInvalidURLParam, op)
		return
	}

	userID := h.authMiddleware.GetUserID(r.Context())
	if userID == 0 {
		response.HandleError(w, r, h.log, errorsAPP.ErrUnauthorized, op)
		return
	}

	navContext, err := h.lessonService.GetLessonNavigationContext(r.Context(), userID, lessonID)
	if err != nil {
		response.HandleError(w, r, h.log, err, op)
		return
	}

	response.OK(w, r, converter.LessonNavigationToDTO(navContext))
}
