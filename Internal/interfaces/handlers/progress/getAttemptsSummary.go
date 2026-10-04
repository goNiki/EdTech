package progress

import (
	"net/http"

	"edtech/internal/dto"
	"edtech/internal/infrastructure/logger"
	response "edtech/internal/interfaces/response"
	errorsAPP "edtech/pkg/errors"
)

// GetAttemptsSummary handles GET /api/v1/lessons/{lesson_id}/attempts/summary
func (h *ProgressHandler) GetAttemptsSummary(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.progress.GetAttemptsSummary"
	log := logger.GetLogger(r.Context(), op)

	userID := h.getUserID(r.Context())
	if userID == 0 {
		response.HandleError(w, r, log, errorsAPP.ErrUnauthorized, op)
		return
	}

	lessonID, err := parseIDParam(r, "lesson_id", "lessonId", "id")
	if err != nil {
		response.HandleError(w, r, log, errorsAPP.ErrInvalidURLParam, op)
		return
	}

	summary, err := h.progressService.GetLessonAttemptsSummary(r.Context(), userID, lessonID)
	if err != nil {
		response.HandleError(w, r, log, err, op)
		return
	}

	var lastAttemptDTO *dto.LessonAttemptItemDTO
	if summary.LastAttempt != nil {
		lastAttemptDTO = &dto.LessonAttemptItemDTO{
			AttemptID:   summary.LastAttempt.AttemptID,
			Score:       summary.LastAttempt.Score,
			SubmittedAt: summary.LastAttempt.SubmittedAt,
		}
	}

	historyDTO := make([]dto.LessonAttemptItemDTO, 0, len(summary.AttemptsHistory))
	for _, hItem := range summary.AttemptsHistory {
		historyDTO = append(historyDTO, dto.LessonAttemptItemDTO{
			AttemptID:   hItem.AttemptID,
			Score:       hItem.Score,
			SubmittedAt: hItem.SubmittedAt,
		})
	}

	response.OK(w, r, dto.LessonAttemptsSummaryResponse{
		LessonID:            summary.LessonID,
		TotalAttemptsMade:   summary.TotalAttemptsMade,
		MaxAttemptsAllowed:  summary.MaxAttemptsAllowed,
		CanStartNewAttempt:  summary.CanStartNewAttempt,
		BestScore:           summary.BestScore,
		BestScorePercentage: summary.BestScorePercentage,
		IsPassed:            summary.IsPassed,
		PassingThreshold:    summary.PassingThreshold,
		LastAttempt:         lastAttemptDTO,
		AttemptsHistory:     historyDTO,
	})
}
