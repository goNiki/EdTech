package courses

import (
	"fmt"
	"net/http"
	"strconv"

	"edtech/internal/dto"
	"edtech/internal/infrastructure/logger"
	"edtech/internal/interfaces/handlers/converter"
	response "edtech/internal/interfaces/response"
	errorsAPP "edtech/pkg/errors"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// UpdateCourse handles PATCH /api/v1/courses/{courseid}
func (h *CourseHandler) UpdateCourse(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.courses.UpdateCourse"

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

	var req dto.UpdateCourseRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		response.HandleError(w, r, log, fmt.Errorf("%w: %w", errorsAPP.ErrDecodeJSON, err), op)
		return
	}

	if err := h.validator.Validate(req); err != nil {
		response.HandleError(w, r, log, fmt.Errorf("%w: %w", errorsAPP.ErrCourseValidation, err), op)
		return
	}

	input := converter.UpdateCourseRequestToDomain(req)

	updatedCourse, err := h.courseService.UpdateCourse(r.Context(), courseID, userID, input)
	if err != nil {
		response.HandleError(w, r, log, err, op)
		return
	}

	response.OK(w, r, converter.CourseToDTO(updatedCourse))
}
