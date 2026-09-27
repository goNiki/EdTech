package converter

import (
	"edtech/internal/domain"
	"edtech/internal/dto"
)

func SectionToDTO(d *domain.Section) dto.Section {
	if d == nil {
		return dto.Section{}
	}
	return dto.Section{
		ID:          d.ID,
		CourseID:    d.CourseID,
		Title:       d.Title,
		Description: d.Description,
		Position:    d.Position,
		Status:      d.Status,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
	}
}

func CreateSectionRequestToDomain(req dto.CreateSectionRequest) domain.Section {
	status := req.Status
	if status == "" {
		status = domain.StatusDraft
	}
	return domain.Section{
		CourseID:    req.CourseID,
		Title:       req.Title,
		Description: req.Description,
		Status:      status,
	}
}
