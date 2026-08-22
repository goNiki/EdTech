package converter

import (
	"edtech/internal/domain"
	"edtech/internal/repository/models"
)

func ResourceToDomain(m *models.Resource) *domain.Resource {
	if m == nil { return nil }
	return &domain.Resource{
		ID:            m.ID,
		LessonID:      m.LessonID,
		CourseID:      m.CourseID,
		Title:         m.Title,
		Description:   m.Description,
		Type:          m.Type,
		Path:          m.Path,
		Mime:          m.Mime,
		Size:          m.Size,
		ExternalURL:   m.ExternalURL,
		Duration:      m.Duration,
		OrderPosition: m.OrderPosition,
		IsRequired:    m.IsRequired,
		CreatedAt:     m.CreatedAt,
		DeletedAt:     m.DeletedAt,
	}
}

func ResourceToEntity(d *domain.Resource) *models.Resource {
	if d == nil { return nil }
	return &models.Resource{
		ID:            d.ID,
		LessonID:      d.LessonID,
		CourseID:      d.CourseID,
		Title:         d.Title,
		Description:   d.Description,
		Type:          d.Type,
		Path:          d.Path,
		Mime:          d.Mime,
		Size:          d.Size,
		ExternalURL:   d.ExternalURL,
		Duration:      d.Duration,
		OrderPosition: d.OrderPosition,
		IsRequired:    d.IsRequired,
		CreatedAt:     d.CreatedAt,
		DeletedAt:     d.DeletedAt,
	}
}
