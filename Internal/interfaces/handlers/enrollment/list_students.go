package enrollment

import (
	"net/http"
	"strconv"

	"edtech/internal/dto"
	"edtech/internal/infrastructure/logger"
	response "edtech/internal/interfaces/response"
	errorsAPP "edtech/pkg/errors"

	"github.com/go-chi/chi/v5"
)

func (h *EnrollmentHandler) ListCourseStudents(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.enrollment.ListCourseStudents"

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

	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")

	page, _ := strconv.ParseInt(pageStr, 10, 64)
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.ParseInt(pageSizeStr, 10, 64)
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	students, total, err := h.enrollmentService.ListCourseStudentsWithProgress(r.Context(), teacherID, courseID, page, pageSize)
	if err != nil {
		response.HandleError(w, r, log, err, op)
		return
	}

	dtos := make([]dto.CourseStudentDTO, 0, len(students))
	for _, s := range students {
		dtos = append(dtos, dto.CourseStudentDTO{
			UserID:             s.UserID,
			FirstName:          s.FirstName,
			LastName:           s.LastName,
			Username:           s.Username,
			Email:              s.Email,
			AvatarURL:          s.AvatarURL,
			Role:               s.Role,
			EnrolledAt:         s.EnrolledAt,
			ProgressPercentage: s.ProgressPercentage,
			CompletedLessons:   s.CompletedLessons,
			TotalLessons:       s.TotalLessons,
			AverageScore:       s.AverageScore,
			HasPending:         s.HasPendingHomeworks,
		})
	}

	resp := dto.PaginatedStudentsResponse{
		Students: dtos,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}

	response.OK(w, r, resp)
}
