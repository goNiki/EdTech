package dto

import "time"

type Lesson struct {
	ID          int64      `json:"id"`
	CourseID    int64      `json:"course_id"`
	SectionID   *int64     `json:"section_id,omitempty"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	CoverURL    string     `json:"cover_url"`
	Content     string     `json:"content"`
	Type        string     `json:"type"` // lecture, exercise, quiz, assignment
	Position    int64      `json:"position"`
	Duration    *int       `json:"duration,omitempty"` // seconds
	IsFree      bool       `json:"is_free"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
}

type CreateLessonRequest struct {
	CourseID    int64  `json:"course_id" validate:"required,gt=0"`
	SectionID   *int64 `json:"section_id,omitempty" validate:"omitempty,gt=0"`
	Title       string `json:"title" validate:"required,min=2,max=255"`
	Description string `json:"description,omitempty" validate:"omitempty,max=2000"`
	CoverURL    string `json:"cover_url,omitempty" validate:"omitempty,http_url,max=500"`
	Content     string `json:"content" validate:"required"`
	Type        string `json:"type" validate:"required,oneof=lecture exercise quiz assignment"`
	Duration    *int   `json:"duration,omitempty" validate:"omitempty,gte=0"`
	IsFree      bool   `json:"is_free"`
}

type UpdateLessonRequest struct {
	SectionID   *int64  `json:"section_id,omitempty" validate:"omitempty,gt=0"`
	Title       *string `json:"title,omitempty" validate:"omitempty,min=2,max=255"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=2000"`
	CoverURL    *string `json:"cover_url,omitempty" validate:"omitempty,http_url,max=500"`
	Content     *string `json:"content,omitempty"`
	Type        *string `json:"type,omitempty" validate:"omitempty,oneof=lecture exercise quiz assignment"`
	Duration    *int    `json:"duration,omitempty" validate:"omitempty,gte=0"`
	IsFree      *bool   `json:"is_free,omitempty"`
}

type CreateLessonResponse struct {
	Lesson Lesson `json:"lesson"`
}
