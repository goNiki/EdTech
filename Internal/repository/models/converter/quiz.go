package converter

import (
	"edtech/internal/domain"
	"edtech/internal/repository/models"
)

func QuizToDomain(m *models.Quiz) *domain.Quiz {
	if m == nil { return nil }
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
	if d == nil { return nil }
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
