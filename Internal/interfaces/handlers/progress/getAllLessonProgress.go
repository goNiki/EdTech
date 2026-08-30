package progress

import (
	"net/http"

	"edtech/internal/dto"
	"edtech/internal/infrastructure/logger"
	"edtech/internal/interfaces/handlers/converter"
	response "edtech/internal/interfaces/response"
	errorsAPP "edtech/pkg/errors"
)

// GetAllLessonProgress handles GET /api/v1/courses/{course_id}/progress/lessons
func (h *ProgressHandler) GetAllLessonProgress(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.progress.GetAllLessonProgress"
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

	progressList, err := h.progressService.GetAllLessonProgress(r.Context(), userID, courseID)
	if err != nil {
		response.HandleError(w, r, log, err, op)
		return
	}

	dtos := make([]dto.LessonProgress, 0, len(progressList))
	for _, p := range progressList {
		dtos = append(dtos, converter.LessonProgressToDTO(&p))
	}

	response.OK(w, r, dtos)
}
