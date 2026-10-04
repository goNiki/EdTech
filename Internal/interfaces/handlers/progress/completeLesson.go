package progress

import (
	"net/http"

	"edtech/internal/domain"
	"edtech/internal/dto"
	"edtech/internal/infrastructure/logger"
	response "edtech/internal/interfaces/response"
	errorsAPP "edtech/pkg/errors"

	"github.com/go-chi/render"
)

// CompleteLesson handles POST /api/v1/lessons/{lesson_id}/complete
func (h *ProgressHandler) CompleteLesson(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.progress.CompleteLesson"
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

	var req dto.CompleteLessonRequest
	if r.Body != nil {
		if err := render.DecodeJSON(r.Body, &req); err != nil && err.Error() != "EOF" {
			log.Warn("failed to decode complete lesson request body", "error", err)
		}
	}

	essays := make([]domain.EssaySubmission, 0, len(req.Essays))
	for _, e := range req.Essays {
		essays = append(essays, domain.EssaySubmission{
			QuestionText: e.QuestionText,
			AnswerText:   e.AnswerText,
			MaxPoints:    e.MaxPoints,
		})
	}

	answers := make([]domain.LessonAnswerSubmission, 0, len(req.Answers))
	for _, a := range req.Answers {
		ansVal := a.Answer
		if ansVal == nil {
			switch {
			case a.SelectedOption != nil:
				ansVal = map[string]any{"selected_option": *a.SelectedOption}
			case len(a.SelectedOpts) > 0:
				ansVal = map[string]any{"selected_options": a.SelectedOpts}
			case a.Pairs != nil:
				ansVal = map[string]any{"pairs": a.Pairs}
			case a.Blanks != nil:
				ansVal = map[string]any{"blanks": a.Blanks}
			case a.Order != nil:
				ansVal = map[string]any{"order": a.Order}
			}
		}
		answers = append(answers, domain.LessonAnswerSubmission{
			BlockID: a.BlockID,
			Answer:  ansVal,
		})
	}

	input := domain.CompleteLessonInput{
		Score:       req.Score,
		TimeSpent:   req.TimeSpent,
		Answers:     answers,
		Essays:      essays,
		AttemptID:   req.AttemptID,
		IsAbandoned: req.IsAbandoned,
	}

	result, err := h.progressService.CompleteLesson(r.Context(), userID, lessonID, input)
	if err != nil {
		response.HandleError(w, r, log, err, op)
		return
	}

	resultsDTO := make(map[string]dto.BlockValidationDTO, len(result.Results))
	for bid, res := range result.Results {
		resultsDTO[bid] = dto.BlockValidationDTO{
			IsCorrect:     res.IsCorrect,
			Feedback:      res.Feedback,
			CorrectAnswer: res.CorrectAnswer,
		}
	}

	response.OK(w, r, dto.CompleteLessonResponse{
		Message:        "lesson completed",
		LessonID:       result.LessonID,
		Status:         result.Status,
		Score:          result.Score,
		EarnedPoints:   result.EarnedPoints,
		TotalPoints:    result.TotalMaxPoints,
		TotalMaxPoints: result.TotalMaxPoints,
		IsPassed:       result.IsPassed,
		Results:        resultsDTO,
	})
}
