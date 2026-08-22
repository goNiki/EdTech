package converter

import (
	"edtech/internal/domain"
	"edtech/internal/dto"
)

func CourseToDTO(d *domain.Course) dto.Course {
	if d == nil {
		return dto.Course{}
	}
	return dto.Course{
		ID:                d.Id,
		Title:             d.Title,
		Slug:              d.Slug,
		ShortDescription:  d.ShortDescription,
		Description:       d.Description,
		CoverURL:          d.CoverURL,
		IntroVideoURL:     d.IntroVideoURL,
		CreatedBy:         d.CreatedBy,
		Visibility:        d.Visibility,
		Status:            d.Status,
		Difficulty:        d.Difficulty,
		Language:          d.Language,
		EstimatedDuration: d.EstimatedDuration,
		CategoryID:        d.CategoryID,
		TotalLessons:      d.TotalLessons,
		TotalSections:     d.TotalSections,
		EnrolledCount:     d.EnrolledCount,
		CreatedAt:         d.CreatedAt,
		UpdatedAt:         d.UpdatedAt,
		PublishedAt:       d.PublishedAt,
	}
}

func CreateCourseRequestToDomain(req dto.CreateCourseRequest, userID int64) domain.Course {
	return domain.Course{
		Title:            req.Title,
		Slug:             req.Slug,
		ShortDescription: req.ShortDescription,
		Description:      req.Description,
		CoverURL:         req.CoverURL,
		CreatedBy:        userID,
		Visibility:       req.Visibility,
		Status:           domain.StatusDraft,
		Difficulty:       req.Difficulty,
		Language:         req.Language,
		CategoryID:       req.CategoryID,
	}
}
