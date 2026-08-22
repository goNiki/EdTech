package converter

import (
	"edtech/internal/domain"
	"edtech/internal/repository/models"
)

func CourseToDomain(m *models.Course) *domain.Course {
	if m == nil {
		return nil
	}
	return &domain.Course{
		Id: m.ID, Title: m.Title, Slug: m.Slug, ShortDescription: m.ShortDescription,
		Description: m.Description, CoverURL: m.CoverURL, IntroVideoURL: m.IntroVideoURL,
		CreatedBy: m.CreatedBy, Visibility: m.Visibility, Status: m.Status,
		Difficulty: m.Difficulty, Language: m.Language, EstimatedDuration: m.EstimatedDuration,
		CategoryID: m.CategoryID, TotalLessons: m.TotalLessons, TotalSections: m.TotalSections,
		EnrolledCount: m.EnrolledCount, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
		PublishedAt: m.PublishedAt, ArchivedAt: m.ArchivedAt, DeletedAt: m.DeletedAt,
	}
}

func CourseToEntity(d *domain.Course) *models.Course {
	if d == nil {
		return nil
	}
	return &models.Course{
		ID: d.Id, Title: d.Title, Slug: d.Slug, ShortDescription: d.ShortDescription,
		Description: d.Description, CoverURL: d.CoverURL, IntroVideoURL: d.IntroVideoURL,
		CreatedBy: d.CreatedBy, Visibility: d.Visibility, Status: d.Status,
		Difficulty: d.Difficulty, Language: d.Language, EstimatedDuration: d.EstimatedDuration,
		CategoryID: d.CategoryID, TotalLessons: d.TotalLessons, TotalSections: d.TotalSections,
		EnrolledCount: d.EnrolledCount, CreatedAt: d.CreatedAt, UpdatedAt: d.UpdatedAt,
		PublishedAt: d.PublishedAt, ArchivedAt: d.ArchivedAt, DeletedAt: d.DeletedAt,
	}
}
