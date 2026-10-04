package analytics

import (
	"fmt"
	"net/http"
	"strconv"

	"edtech/internal/infrastructure/logger"
	response "edtech/internal/interfaces/response"
	errorsAPP "edtech/pkg/errors"

	"github.com/go-chi/chi/v5"
)

func (h *AnalyticsHandler) ExportGradebook(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.analytics.ExportGradebook"

	log := logger.GetLogger(r.Context(), op)

	courseIDStr := chi.URLParam(r, "courseid")
	if courseIDStr == "" {
		courseIDStr = chi.URLParam(r, "id")
	}

	courseID, err := strconv.ParseInt(courseIDStr, 10, 64)
	if err != nil {
		response.HandleError(w, r, log, errorsAPP.ErrInvalidURLParam, op)
		return
	}

	format := r.URL.Query().Get("format")
	if format != "" && format != "csv" {
		response.HandleError(w, r, log, errorsAPP.ErrInvalidURLQuery, op)
		return
	}

	teacherID := h.authMiddleware.GetUserID(r.Context())
	if teacherID == 0 {
		response.HandleError(w, r, log, errorsAPP.ErrUnauthorized, op)
		return
	}

	data, filename, err := h.analyticsService.ExportCourseGradebookCSV(r.Context(), teacherID, courseID)
	if err != nil {
		response.HandleError(w, r, log, err, op)
		return
	}

	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filename))
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}
