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
