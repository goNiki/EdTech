package courses

import (
	"edtech/internal/infrastructure/logger"
	"edtech/internal/infrastructure/response"
	errorsAPP "edtech/pkg/errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (h *handler) PublishCourse(w http.ResponseWriter, r *http.Request) {
	const op = "http.handler.course.publishcource"

	log := logger.GetLogger(r.Context(), op)

	strcourseID := chi.URLParam(r, "courseid")

	courseID, err := strconv.Atoi(strcourseID)
	if err != nil {
		response.HandleError(w, r, log, errorsAPP.ErrInvalidURLParam, op)
		return
	}
	userID := h.authMiddleware.GetUserID(r.Context())

	if err := h.courseService.PublishCourse(r.Context(), userID, int64(courseID)); err != nil {
		response.HandleError(w, r, log, err, op)
		return
	}

	response.NoContent(w, r)

}
