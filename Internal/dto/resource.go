package dto

import "time"

type Resource struct {
	ID            int64     `json:"id"`
	LessonID      *int64    `json:"lesson_id,omitempty"`
	CourseID      *int64    `json:"course_id,omitempty"`
	Title         *string   `json:"title,omitempty"`
	Description   *string   `json:"description,omitempty"`
	Type          string    `json:"type"`
	Path          string    `json:"path"`
	Mime          string    `json:"mime"`
	Size          int64     `json:"size"`
	ExternalURL   *string   `json:"external_url,omitempty"`
	Duration      *int      `json:"duration,omitempty"`
	OrderPosition int       `json:"order_position"`
	IsRequired    bool      `json:"is_required"`
	CreatedAt     time.Time `json:"created_at"`
}

type UploadResourceRequest struct {
	LessonID    *int64  `json:"lesson_id,omitempty" validate:"omitempty,gt=0"`
	CourseID    *int64  `json:"course_id,omitempty" validate:"omitempty,gt=0"`
	Title       *string `json:"title,omitempty" validate:"omitempty,min=1,max=255"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=1000"`
	ExternalURL *string `json:"external_url,omitempty" validate:"omitempty,http_url,max=1000"`
	Duration    *int    `json:"duration,omitempty" validate:"omitempty,gte=0"`
	IsRequired  bool    `json:"is_required"`
}

type ResourceResponse struct {
	Resource Resource `json:"resource"`
}
