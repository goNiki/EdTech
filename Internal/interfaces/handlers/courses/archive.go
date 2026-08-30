package courses

import (
	"net/http"
	"strconv"

	"edtech/internal/infrastructure/logger"
	response "edtech/internal/interfaces/response"
	errorsAPP "edtech/pkg/errors"

	"github.com/go-chi/chi/v5"
)

// ArchiveCourse handles POST /api/v1/courses/{courseid}/archive
func (h *CourseHandler) ArchiveCourse(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.courses.ArchiveCourse"

	log := logger.GetLogger(r.Context(), op)

	courseIDStr := chi.URLParam(r, "courseid")
	courseID, err := strconv.ParseInt(courseIDStr, 10, 64)
	if err != nil {
		response.HandleError(w, r, log, errorsAPP.ErrInvalidURLParam, op)
		return
	}

	userID := h.authMiddleware.GetUserID(r.Context())
	if userID == 0 {
		response.HandleError(w, r, log, errorsAPP.ErrUnauthorized, op)
		return
	}

	if err := h.courseService.ArchiveCourse(r.Context(), userID, courseID); err != nil {
		response.HandleError(w, r, log, err, op)
		return
	}

	response.NoContent(w, r)
}
