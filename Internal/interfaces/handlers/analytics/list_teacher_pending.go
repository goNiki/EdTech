package analytics

import (
	"net/http"
	"strconv"

	"edtech/internal/domain"
	"edtech/internal/dto"
	"edtech/internal/infrastructure/logger"
	response "edtech/internal/interfaces/response"
	errorsAPP "edtech/pkg/errors"
)

func (h *AnalyticsHandler) ListTeacherPendingHomeworks(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.analytics.ListTeacherPendingHomeworks"

	log := logger.GetLogger(r.Context(), op)

	teacherID := h.authMiddleware.GetUserID(r.Context())
	if teacherID == 0 {
		response.HandleError(w, r, log, errorsAPP.ErrUnauthorized, op)
		return
	}

	userRole := h.authMiddleware.GetUserRole(r.Context())
	if userRole != string(domain.RoleTeacher) && userRole != string(domain.RoleAdmin) {
		response.HandleError(w, r, log, errorsAPP.ErrForbidden, op)
		return
	}

	var courseID int64
	if courseIDStr := r.URL.Query().Get("course_id"); courseIDStr != "" {
		if parsed, err := strconv.ParseInt(courseIDStr, 10, 64); err == nil && parsed > 0 {
			courseID = parsed
		}
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

	res, err := h.analyticsService.ListTeacherPendingHomeworks(r.Context(), teacherID, courseID, page, pageSize)
	if err != nil {
		response.HandleError(w, r, log, err, op)
		return
	}

	itemsDTO := make([]dto.PendingHomeworkDTO, 0, len(res.Items))
	for _, it := range res.Items {
		itemsDTO = append(itemsDTO, dto.PendingHomeworkDTO{
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

	summaryDTO := make([]dto.CoursePendingSummaryDTO, 0, len(res.CoursesSummary))
	for _, s := range res.CoursesSummary {
		summaryDTO = append(summaryDTO, dto.CoursePendingSummaryDTO{
			CourseID:     s.CourseID,
			CourseTitle:  s.CourseTitle,
			PendingCount: s.PendingCount,
		})
	}

	resp := dto.PaginatedTeacherPendingHomeworksResponse{
		Items:          itemsDTO,
		Total:          res.Total,
		Page:           page,
		PageSize:       pageSize,
		CoursesSummary: summaryDTO,
	}

	response.OK(w, r, resp)
}
