package enrollment

import (
	"net/http"
	"strconv"

	"edtech/internal/infrastructure/logger"
	response "edtech/internal/interfaces/response"
	errorsAPP "edtech/pkg/errors"

	"github.com/go-chi/chi/v5"
)

func (h *EnrollmentHandler) UnenrollSelf(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.enrollment.UnenrollSelf"

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

	if err := h.enrollmentService.UnenrollUser(r.Context(), userID, courseID); err != nil {
		response.HandleError(w, r, log, err, op)
		return
	}

	response.OK(w, r, map[string]interface{}{
		"success": true,
		"message": "successfully unenrolled from course",
	})
}

func (h *EnrollmentHandler) TeacherUnenroll(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.enrollment.TeacherUnenroll"

	log := logger.GetLogger(r.Context(), op)

	courseIDStr := chi.URLParam(r, "courseid")
	courseID, err := strconv.ParseInt(courseIDStr, 10, 64)
	if err != nil {
		response.HandleError(w, r, log, errorsAPP.ErrInvalidURLParam, op)
		return
	}

	targetUserIDStr := chi.URLParam(r, "userid")
	targetUserID, err := strconv.ParseInt(targetUserIDStr, 10, 64)
	if err != nil {
		response.HandleError(w, r, log, errorsAPP.ErrInvalidURLParam, op)
		return
	}

	if err := h.enrollmentService.UnenrollUser(r.Context(), targetUserID, courseID); err != nil {
		response.HandleError(w, r, log, err, op)
		return
	}

	response.OK(w, r, map[string]interface{}{
		"success": true,
		"message": "student successfully removed from course",
	})
}
