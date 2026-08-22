package dto

import "time"

type Section struct {
	ID          int64      `json:"id"`
	CourseID    int64      `json:"course_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Position    int        `json:"position"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type CreateSectionRequest struct {
	CourseID    int64  `json:"course_id" validate:"required,gt=0"`
	Title       string `json:"title" validate:"required,min=2,max=255"`
	Description string `json:"description,omitempty" validate:"omitempty,max=1000"`
}

type UpdateSectionRequest struct {
	Title       *string `json:"title,omitempty" validate:"omitempty,min=2,max=255"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=1000"`
}

type SectionWithLessons struct {
	Section Section  `json:"section"`
	Lessons []Lesson `json:"lessons"`
}
