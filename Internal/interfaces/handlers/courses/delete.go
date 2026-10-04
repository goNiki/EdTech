package courses

import (
	"edtech/internal/infrastructure/logger"
	response "edtech/internal/interfaces/response"
	errorsAPP "edtech/pkg/errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (h *CourseHandler) DeleteCourse(w http.ResponseWriter, r *http.Request) {
	const op = "http.handler.course.deletecourse"

	log := logger.GetLogger(r.Context(), op)

	strcourseID := chi.URLParam(r, "courseid")

	courseID, err := strconv.Atoi(strcourseID)
	if err != nil {
		response.HandleError(w, r, log, errorsAPP.ErrInvalidURLParam, op)
		return
	}

	userID := h.authMiddleware.GetUserID(r.Context())
	if userID == 0 {
		response.HandleError(w, r, log, errorsAPP.ErrUnauthorized, op)
		return
	}

	if err := h.courseService.DeleteCourse(r.Context(), int64(courseID), userID); err != nil {
		response.HandleError(w, r, log, err, op)
		return
	}

	response.OK(w, r, map[string]string{"message": "course successfully deleted"})
}
