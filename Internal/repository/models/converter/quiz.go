package converter

import (
	"edtech/internal/domain"
	"edtech/internal/repository/models"
)

func QuizToDomain(m *models.Quiz) *domain.Quiz {
	if m == nil {
		return nil
	}
	return &domain.Quiz{
		ID:          m.ID,
		LessonID:    m.LessonID,
		Title:       m.Title,
		Description: m.Description,
		PassingScor: m.PassingScor,
		MaxAttempts: m.MaxAttempts,
		TimeLimit:   m.TimeLimit,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		DeletedAt:   m.DeletedAt,
	}
}

func QuizToEntity(d *domain.Quiz) *models.Quiz {
	if d == nil {
		return nil
	}
	return &models.Quiz{
		ID:          d.ID,
		LessonID:    d.LessonID,
		Title:       d.Title,
		Description: d.Description,
		PassingScor: d.PassingScor,
		MaxAttempts: d.MaxAttempts,
		TimeLimit:   d.TimeLimit,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
		DeletedAt:   d.DeletedAt,
	}
}

func QuizQuestionToDomain(m *models.QuizQuestion) *domain.QuizQuestion {
	if m == nil {
		return nil
	}
	return &domain.QuizQuestion{
		ID:       m.ID,
		QuizID:   m.QuizID,
		Type:     domain.QuizType(m.Type),
		Text:     m.Text,
		Points:   m.Points,
		Position: m.Position,
	}
}

func QuizQuestionToEntity(d *domain.QuizQuestion) *models.QuizQuestion {
	if d == nil {
		return nil
	}
	return &models.QuizQuestion{
		ID:       d.ID,
		QuizID:   d.QuizID,
		Type:     string(d.Type),
		Text:     d.Text,
		Points:   d.Points,
		Position: d.Position,
	}
}

func QuizAnswerToDomain(m *models.QuizAnswer) *domain.QuizAnswer {
	if m == nil {
		return nil
	}
	return &domain.QuizAnswer{
		ID:         m.ID,
		QuestionID: m.QuestionID,
		Text:       m.Text,
		IsCorrect:  m.IsCorrect,
		Explain:    m.Explain,
	}
}

func QuizAnswerToEntity(d *domain.QuizAnswer) *models.QuizAnswer {
	if d == nil {
		return nil
	}
	return &models.QuizAnswer{
		ID:         d.ID,
		QuestionID: d.QuestionID,
		Text:       d.Text,
		IsCorrect:  d.IsCorrect,
		Explain:    d.Explain,
	}
}

func QuizAttemptToDomain(m *models.QuizAttempt) *domain.QuizAttempt {
	if m == nil {
		return nil
	}
	return &domain.QuizAttempt{
		ID:           m.ID,
		QuizID:       m.QuizID,
		UserID:       m.UserID,
		Score:        m.Score,
		Passed:       m.Passed,
		NeedsGrading: m.NeedsGrading,
		StartedAt:    m.StartedAt,
		CompletedAt:  m.CompletedAt,
	}
}

func QuizAttemptToEntity(d *domain.QuizAttempt) *models.QuizAttempt {
	if d == nil {
		return nil
	}
	return &models.QuizAttempt{
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

func QuizAttemptAnswerToDomain(m *models.QuizAttemptAnswer) *domain.QuizAttemptAnswer {
	if m == nil {
		return nil
	}
	return &domain.QuizAttemptAnswer{
		ID:         m.ID,
		AttemptID:  m.AttemptID,
		QuestionID: m.QuestionID,
		AnswerID:   m.AnswerID,
		TextValue:  m.TextValue,
		IsCorrect:  m.IsCorrect,
		Points:     m.Points,
		Feedback:   m.Feedback,
	}
}

func QuizAttemptAnswerToEntity(d *domain.QuizAttemptAnswer) *models.QuizAttemptAnswer {
	if d == nil {
		return nil
	}
	return &models.QuizAttemptAnswer{
		ID:         d.ID,
		AttemptID:  d.AttemptID,
		QuestionID: d.QuestionID,
		AnswerID:   d.AnswerID,
		TextValue:  d.TextValue,
		IsCorrect:  d.IsCorrect,
		Points:     d.Points,
		Feedback:   d.Feedback,
	}
}
