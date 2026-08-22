package converter

import (
	"edtech/internal/domain"
	"edtech/internal/repository/models"
)

func LessonToDomain(m *models.Lesson) *domain.Lesson {
	if m == nil { return nil }
	return &domain.Lesson{
		ID:          m.ID,
		CourseID:    m.CourseID,
		SectionID:   m.SectionID,
		Title:       m.Title,
		Description: m.Description,
		CoverURL:    m.CoverURL,
		Content:     m.Content,
		Type:        m.Type,
		Position:    m.Position,
		Duration:    m.Duration,
		IsFree:      m.IsFree,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		PublishedAt: m.PublishedAt,
		DeletedAt:   m.DeletedAt,
	}
}

func LessonToEntity(d *domain.Lesson) *models.Lesson {
	if d == nil { return nil }
	return &models.Lesson{
		ID:          d.ID,
		CourseID:    d.CourseID,
		SectionID:   d.SectionID,
		Title:       d.Title,
		Description: d.Description,
		CoverURL:    d.CoverURL,
		Content:     d.Content,
		Type:        d.Type,
		Position:    d.Position,
		Duration:    d.Duration,
		IsFree:      d.IsFree,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
		PublishedAt: d.PublishedAt,
		DeletedAt:   d.DeletedAt,
	}
}
