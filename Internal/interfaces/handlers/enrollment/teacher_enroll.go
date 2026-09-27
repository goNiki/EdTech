package enrollment

import (
	"net/http"
	"strconv"

	"edtech/internal/domain"
	"edtech/internal/dto"
	"edtech/internal/infrastructure/logger"
	response "edtech/internal/interfaces/response"
	errorsAPP "edtech/pkg/errors"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func (h *EnrollmentHandler) TeacherEnroll(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.enrollment.TeacherEnroll"

	log := logger.GetLogger(r.Context(), op)

	courseIDStr := chi.URLParam(r, "courseid")
	courseID, err := strconv.ParseInt(courseIDStr, 10, 64)
	if err != nil {
		response.HandleError(w, r, log, errorsAPP.ErrInvalidURLParam, op)
		return
	}

	teacherID := h.authMiddleware.GetUserID(r.Context())
	if teacherID == 0 {
		response.HandleError(w, r, log, errorsAPP.ErrUnauthorized, op)
		return
	}

	var req dto.TeacherEnrollRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		response.HandleError(w, r, log, errorsAPP.ErrDecodeJSON, op)
		return
	}

	domainReq := domain.TeacherEnrollRequest{
		TeacherID:    teacherID,
		CourseID:     courseID,
		Role:         req.Role,
		TargetEmail:  req.Email,
		TargetUserID: req.UserID,
	}

	if err := h.enrollmentService.TeacherEnrollCourse(r.Context(), domainReq); err != nil {
		response.HandleError(w, r, log, err, op)
		return
	}

	response.OK(w, r, map[string]interface{}{
		"success": true,
		"message": "student successfully added to course",
	})
}
