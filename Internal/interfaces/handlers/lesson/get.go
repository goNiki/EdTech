package lesson

import (
	"edtech/internal/interfaces/handlers/converter"
	"edtech/internal/interfaces/response"
	errorsAPP "edtech/pkg/errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (h *LessonHandler) GetLesson(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.lesson.GetLesson"

	idStr := chi.URLParam(r, "id")
	lessonID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.HandleError(w, r, h.log, errorsAPP.ErrInvalidURLParam, op)
		return
	}

	lesson, err := h.lessonService.GetLesson(r.Context(), lessonID)
	if err != nil {
		response.HandleError(w, r, h.log, err, op)
		return
	}

	response.OK(w, r, converter.LessonToDTO(lesson))
}
