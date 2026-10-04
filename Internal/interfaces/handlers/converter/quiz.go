package converter

import (
	"edtech/internal/domain"
	"edtech/internal/dto"
)

func QuizToDTO(d *domain.Quiz) dto.Quiz {
	if d == nil {
		return dto.Quiz{}
	}
	return dto.Quiz{
		ID:          d.ID,
		LessonID:    d.LessonID,
		Title:       d.Title,
		Description: d.Description,
		PassingScor: d.PassingScor,
		MaxAttempts: d.MaxAttempts,
		TimeLimit:   d.TimeLimit,
		CreatedAt:   d.CreatedAt,
	}
}

func QuizQuestionToDTO(d *domain.QuizQuestion) dto.QuizQuestion {
	if d == nil {
		return dto.QuizQuestion{}
	}
	return dto.QuizQuestion{
		ID:       d.ID,
		QuizID:   d.QuizID,
		Type:     string(d.Type),
		Text:     d.Text,
		Points:   d.Points,
		Position: d.Position,
	}
}

func QuizAttemptToDTO(d *domain.QuizAttempt) dto.QuizAttempt {
	if d == nil {
		return dto.QuizAttempt{}
	}
	return dto.QuizAttempt{
		ID:           d.ID,
		QuizID:       d.QuizID,
		UserID:       d.UserID,
		Score:        d.Score,
		Passed:       d.Passed,
		NeedsGrading: d.NeedsGrading,
		StartedAt:    d.StartedAt,
		CompletedAt:  d.CompletedAt,
	}
}

func QuizAttemptsToDTO(attempts []domain.QuizAttempt) []dto.QuizAttempt {
	dtos := make([]dto.QuizAttempt, 0, len(attempts))
	for _, a := range attempts {
		dtos = append(dtos, QuizAttemptToDTO(&a))
	}
	return dtos
}

func StudentHomeworkFeedbackToDTO(f *domain.StudentHomeworkFeedback) dto.HomeworkFeedbackResponse {
	if f == nil {
		return dto.HomeworkFeedbackResponse{
			Code:    200,
			Message: "success",
			Data: dto.HomeworkFeedbackDataDTO{
				HasSubmission: false,
				Status:        "not_submitted",
				Answers:       make([]dto.HomeworkAnswerDetailDTO, 0),
			},
		}
	}

	var teacherDTO *dto.HomeworkTeacherDTO
	if f.Teacher != nil {
		teacherDTO = &dto.HomeworkTeacherDTO{
			ID:        f.Teacher.ID,
			Name:      f.Teacher.Name,
			AvatarURL: f.Teacher.AvatarURL,
		}
	}

	answersDTO := make([]dto.HomeworkAnswerDetailDTO, 0, len(f.Answers))
	for _, a := range f.Answers {
		answersDTO = append(answersDTO, dto.HomeworkAnswerDetailDTO{
			AnswerID:      a.AnswerID,
			QuestionText:  a.QuestionText,
			StudentAnswer: a.StudentAnswer,
			Points:        a.Points,
			MaxPoints:     a.MaxPoints,
			IsCorrect:     a.IsCorrect,
			Feedback:      a.Feedback,
		})
	}

	return dto.HomeworkFeedbackResponse{
		Code:    200,
		Message: "success",
		Data: dto.HomeworkFeedbackDataDTO{
			HasSubmission: f.HasSubmission,
			Status:        f.Status,
			AttemptID:     f.AttemptID,
			SubmittedAt:   f.SubmittedAt,
			GradedAt:      f.GradedAt,
			Teacher:       teacherDTO,
			Answers:       answersDTO,
		},
	}
}

