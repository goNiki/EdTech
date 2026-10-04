package converter

import (
	"edtech/internal/domain"
	"edtech/internal/dto"
)

func CategoryToDTO(d *domain.Category) dto.Category {
	if d == nil {
		return dto.Category{}
	}
	return dto.Category{
		ID:           d.ID,
		Name:         d.Name,
		Slug:         d.Slug,
		Description:  d.Description,
		IconURL:      d.IconURL,
		ParentID:     d.ParentID,
		CoursesCount: d.CoursesCount,
	}
}

func CategoriesToDTO(list []domain.Category) []dto.Category {
	res := make([]dto.Category, 0, len(list))
	for i := range list {
		res = append(res, CategoryToDTO(&list[i]))
	}
	return res
}

func CreateCategoryRequestToDomain(req dto.CreateCategoryRequest) domain.Category {
	return domain.Category{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		ParentID:    req.ParentID,
	}
}
