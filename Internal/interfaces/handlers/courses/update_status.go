package courses

import (
	"net/http"
	"strconv"

	"edtech/internal/dto"
	"edtech/internal/infrastructure/logger"
	response "edtech/internal/interfaces/response"
	errorsAPP "edtech/pkg/errors"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// UpdateCourseStatus handles PATCH /api/v1/courses/{courseid}/status
func (h *CourseHandler) UpdateCourseStatus(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.courses.UpdateCourseStatus"

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

	var req dto.UpdateStatusRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		response.HandleError(w, r, log, errorsAPP.ErrDecodeJSON, op)
		return
	}

	if err := h.validator.Validate(req); err != nil {
		response.HandleError(w, r, log, errorsAPP.ErrCourseValidation, op)
		return
	}

	if err := h.courseService.UpdateCourseStatus(r.Context(), userID, courseID, req.Status); err != nil {
		response.HandleError(w, r, log, err, op)
		return
	}

	response.OK(w, r, map[string]interface{}{
		"success": true,
		"message": "course status updated",
		"status":  req.Status,
	})
}
