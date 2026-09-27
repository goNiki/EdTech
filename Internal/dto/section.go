package dto

import "time"

type Section struct {
	ID          int64     `json:"id"`
	CourseID    int64     `json:"course_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Position    int       `json:"position"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateSectionRequest struct {
	CourseID    int64  `json:"course_id" validate:"required,gt=0"`
	Title       string `json:"title" validate:"required,min=1,max=255"`
	Description string `json:"description,omitempty" validate:"omitempty,max=1000"`
	Status      string `json:"status,omitempty" validate:"omitempty,oneof=draft published archived"`
}

type UpdateSectionRequest struct {
	Title       *string `json:"title,omitempty" validate:"omitempty,min=1,max=255"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=1000"`
	Position    *int    `json:"position,omitempty"`
	Status      *string `json:"status,omitempty" validate:"omitempty,oneof=draft published archived"`
}

type SectionWithLessons struct {
	Section Section  `json:"section"`
	Lessons []Lesson `json:"lessons"`
}
