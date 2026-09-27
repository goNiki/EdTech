package converter

import (
	"edtech/internal/domain"
	"edtech/internal/dto"
)

func LessonToDTO(d *domain.Lesson) dto.Lesson {
	if d == nil {
		return dto.Lesson{}
	}
	return dto.Lesson{
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
		Status:      d.Status,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
		PublishedAt: d.PublishedAt,
	}
}

func CreateLessonRequestToDomain(req dto.CreateLessonRequest) domain.Lesson {
	status := req.Status
	if status == "" {
		status = domain.StatusDraft
	}
	return domain.Lesson{
		CourseID:    req.CourseID,
		SectionID:   req.SectionID,
		Title:       req.Title,
		Description: req.Description,
		CoverURL:    req.CoverURL,
		Content:     req.Content,
		Type:        req.Type,
		Duration:    req.Duration,
		IsFree:      req.IsFree,
		Status:      status,
	}
}
