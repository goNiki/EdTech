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
		Rating:            d.Rating,
		ReviewsCount:      d.ReviewsCount,
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

func UpdateCourseRequestToDomain(req dto.UpdateCourseRequest) domain.UpdateCourseInput {
	return domain.UpdateCourseInput{
		Title:             req.Title,
		Slug:              req.Slug,
		ShortDescription:  req.ShortDescription,
		Description:       req.Description,
		CoverURL:          req.CoverURL,
		IntroVideoURL:     req.IntroVideoURL,
		Visibility:        req.Visibility,
		Difficulty:        req.Difficulty,
		Language:          req.Language,
		EstimatedDuration: req.EstimatedDuration,
		CategoryID:        req.CategoryID,
	}
}

func ListPublicCoursesRequestToDomain(req dto.ListPublicCoursesRequest) (domain.Pagination, domain.CourseFilter) {
	pagination := domain.Pagination{
		Page:     req.Pagination.Page,
		PageSize: req.Pagination.PageSize,
	}
	filter := domain.CourseFilter{
		Search:     req.Filter.Search,
		CategoryID: req.Filter.CategoryID,
		CreatedBy:  req.Filter.CreatedBy,
		Difficulty: req.Filter.Difficulty,
		Language:   req.Filter.Language,
		SortBy:     req.Filter.SortBy,
		SortOrder:  req.Filter.SortOrder,
	}
	return pagination, filter
}

func ListMyCoursesRequestToDomain(req dto.ListMyCoursesRequest) *domain.InputListMyCourse {
	return &domain.InputListMyCourse{
		UserID: req.UserID,
		Role:   req.Role,
		Pagination: domain.Pagination{
			Page:     req.Pagination.Page,
			PageSize: req.Pagination.PageSize,
		},
		Filter: domain.CourseFilter{
			Search:     req.Filter.Search,
			CategoryID: req.Filter.CategoryID,
			CreatedBy:  req.Filter.CreatedBy,
			Difficulty: req.Filter.Difficulty,
			Language:   req.Filter.Language,
			SortBy:     req.Filter.SortBy,
			SortOrder:  req.Filter.SortOrder,
		},
	}
}

func CoursePermissionsToDTO(p domain.CoursePermissions) dto.CoursePermissions {
	return dto.CoursePermissions{
		CanView:        p.CanView,
		CanEdit:        p.CanEdit,
		CanDelete:      p.CanDelete,
		CanPublish:     p.CanPublish,
		CanEnroll:      p.CanEnroll,
		CanManageUsers: p.CanManageUsers,
	}
}
