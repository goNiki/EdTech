package analytics

import (
	"net/http"
	"strconv"

	"edtech/internal/dto"
	"edtech/internal/infrastructure/logger"
	response "edtech/internal/interfaces/response"
	errorsAPP "edtech/pkg/errors"

	"github.com/go-chi/chi/v5"
)

func (h *AnalyticsHandler) GetCourseAnalytics(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.analytics.GetCourseAnalytics"

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

	summary, err := h.analyticsService.GetCourseAnalytics(r.Context(), teacherID, courseID)
	if err != nil {
		response.HandleError(w, r, log, err, op)
		return
	}

	resp := dto.CourseAnalyticsResponse{
		TotalStudents:         summary.TotalStudents,
		AvgProgressPercent:    summary.AvgProgressPercent,
		AvgScore:              summary.AvgScore,
		PendingHomeworksCount: summary.PendingHomeworksCount,
	}

	response.OK(w, r, resp)
}
