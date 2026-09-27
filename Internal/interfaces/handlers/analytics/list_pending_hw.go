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

func (h *AnalyticsHandler) ListPendingHomeworks(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.analytics.ListPendingHomeworks"

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

	items, total, err := h.analyticsService.ListPendingHomeworks(r.Context(), teacherID, courseID, page, pageSize)
	if err != nil {
		response.HandleError(w, r, log, err, op)
		return
	}

	dtos := make([]dto.PendingHomeworkDTO, 0, len(items))
	for _, it := range items {
		dtos = append(dtos, dto.PendingHomeworkDTO{
			AttemptID:       it.AttemptID,
			AnswerID:        it.AnswerID,
			StudentID:       it.StudentID,
			StudentName:     it.StudentName,
			StudentEmail:    it.StudentEmail,
			StudentUsername: it.StudentUsername,
			CourseID:        it.CourseID,
			CourseTitle:     it.CourseTitle,
			LessonID:        it.LessonID,
			LessonTitle:     it.LessonTitle,
			QuestionID:      it.QuestionID,
			QuestionText:    it.QuestionText,
			StudentAnswer:   it.StudentAnswer,
			AttachmentURL:   it.AttachmentURL,
			MaxPoints:       it.MaxPoints,
			Rubric:          it.Rubric,
			SubmittedAt:     it.SubmittedAt,
		})
	}

	resp := dto.PaginatedPendingHomeworksResponse{
		Items:    dtos,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}

	response.OK(w, r, resp)
}
