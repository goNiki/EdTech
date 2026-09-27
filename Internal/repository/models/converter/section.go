package converter

import (
	"edtech/internal/domain"
	"edtech/internal/repository/models"
)

func SectionToDomain(m *models.Section) *domain.Section {
	if m == nil {
		return nil
	}
	return &domain.Section{
		ID:          m.ID,
		CourseID:    m.CourseID,
		Title:       m.Title,
		Description: m.Description,
		Position:    m.Position,
		Status:      m.Status,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		DeletedAt:   m.DeletedAt,
	}
}

func SectionToEntity(d *domain.Section) *models.Section {
	if d == nil {
		return nil
	}
	return &models.Section{
		ID:          d.ID,
		CourseID:    d.CourseID,
		Title:       d.Title,
		Description: d.Description,
		Position:    d.Position,
		Status:      d.Status,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
		DeletedAt:   d.DeletedAt,
	}
}
