package progress

import (
	"net/http"

	"edtech/internal/infrastructure/logger"
	"edtech/internal/interfaces/handlers/converter"
	response "edtech/internal/interfaces/response"
	errorsAPP "edtech/pkg/errors"
)

// GetCourseProgress handles GET /api/v1/courses/{course_id}/progress
func (h *ProgressHandler) GetCourseProgress(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.progress.GetCourseProgress"
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

	prog, err := h.progressService.GetCourseProgress(r.Context(), userID, courseID)
	if err != nil {
		response.HandleError(w, r, log, err, op)
		return
	}

	response.OK(w, r, converter.CourseProgressToDTO(prog))
}
