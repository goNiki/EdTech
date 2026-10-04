package dto

type Category struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	Description  string `json:"description"`
	IconURL      string `json:"icon_url,omitempty"`
	ParentID     *int64 `json:"parent_id,omitempty"`
	CoursesCount int    `json:"courses_count"`
}

type CategoriesResponse struct {
	Categories []Category `json:"categories"`
}

type CreateCategoryRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=100"`
	Slug        string `json:"slug" validate:"required,min=2,max=100"`
	Description string `json:"description,omitempty" validate:"omitempty,max=500"`
	ParentID    *int64 `json:"parent_id,omitempty" validate:"omitempty,gt=0"`
}

type UpdateCategoryRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Slug        *string `json:"slug,omitempty" validate:"omitempty,min=2,max=100"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=500"`
	ParentID    *int64  `json:"parent_id,omitempty" validate:"omitempty,gt=0"`
}
