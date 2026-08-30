package converter

import (
	"edtech/internal/domain"
	"edtech/internal/repository/models"
)

func CategoryToDomain(m *models.Category) *domain.Category {
	if m == nil {
		return nil
	}
	return &domain.Category{
		ID:          m.ID,
		Name:        m.Name,
		Slug:        m.Slug,
		Description: m.Description,
		ParentID:    m.ParentID,
	}
}

func CategoryToEntity(d *domain.Category) *models.Category {
	if d == nil {
		return nil
	}
	return &models.Category{
		ID:          d.ID,
		Name:        d.Name,
		Slug:        d.Slug,
		Description: d.Description,
		ParentID:    d.ParentID,
	}
}
