package converter

import (
	"edtech/internal/domain"
	"edtech/internal/repository/models"
)

func CourseToDomain(m *models.Course) *domain.Course {
	if m == nil {
		return nil
	}
	res := &domain.Course{}
	res.Id = m.ID
	res.Title = m.Title
	res.Slug = m.Slug
	res.ShortDescription = m.ShortDescription
	res.Description = m.Description
	res.CoverURL = m.CoverURL
	res.IntroVideoURL = m.IntroVideoURL
	res.CreatedBy = m.CreatedBy
	res.Visibility = m.Visibility
	res.Status = m.Status
	res.Difficulty = m.Difficulty
	res.Language = m.Language
	res.EstimatedDuration = m.EstimatedDuration
	res.CategoryID = m.CategoryID
	res.TotalLessons = m.TotalLessons
	res.TotalSections = m.TotalSections
	res.EnrolledCount = m.EnrolledCount
	res.Rating = m.Rating
	res.ReviewsCount = m.ReviewsCount
	res.CreatedAt = m.CreatedAt
	res.UpdatedAt = m.UpdatedAt
	res.PublishedAt = m.PublishedAt
	res.ArchivedAt = m.ArchivedAt
	res.DeletedAt = m.DeletedAt
	return res
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
		EnrolledCount: d.EnrolledCount, Rating: d.Rating, ReviewsCount: d.ReviewsCount,
		CreatedAt: d.CreatedAt, UpdatedAt: d.UpdatedAt,
		PublishedAt: d.PublishedAt, ArchivedAt: d.ArchivedAt, DeletedAt: d.DeletedAt,
	}
}
