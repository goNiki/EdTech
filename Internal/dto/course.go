package dto

import (
	"time"
)

const (
	VisibilityPublic  = "public"
	VisibilityPrivate = "private"
)

const (
	StatusDraft     = "draft"
	StatusPublished = "published"
	StatusArchived  = "archived"
)

type Course struct {
	ID                int64      `json:"id"`
	Title             string     `json:"title"`
	Slug              string     `json:"slug"`
	ShortDescription  *string    `json:"short_description,omitempty"`
	Description       string     `json:"description"`
	CoverURL          string     `json:"cover_url"`
	IntroVideoURL     *string    `json:"intro_video_url,omitempty"`
	CreatedBy         int64      `json:"created_by"`
	Visibility        string     `json:"visibility" validate:"required,oneof=public private"`
	Status            string     `json:"status" validate:"required,oneof=draft published archived"`
	Difficulty        *string    `json:"difficulty,omitempty"`
	Language          *string    `json:"language,omitempty"`
	EstimatedDuration *int       `json:"estimated_duration,omitempty"`
	CategoryID        *int64     `json:"category_id,omitempty"`
	TotalLessons      int        `json:"total_lessons"`
	TotalSections     int        `json:"total_sections"`
	EnrolledCount     int        `json:"enrolled_count"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	PublishedAt       *time.Time `json:"published_at,omitempty"`
}

type CoursePermissions struct {
	CanView        bool `json:"can_view"`
	CanEdit        bool `json:"can_edit"`
	CanDelete      bool `json:"can_delete"`
	CanPublish     bool `json:"can_publish"`
	CanEnroll      bool `json:"can_enroll"`
	CanManageUsers bool `json:"can_manage_users"`
}

type CourseDetailResponse struct {
	Course      Course            `json:"course"`
	Permissions CoursePermissions `json:"permissions"`
	UserRole    string            `json:"user_role,omitempty"`
}

type CourseWithPermissions struct {
	Course     Course            `json:"course"`
	Permission CoursePermissions `json:"permissions"`
}

type CreateCourseRequest struct {
	Title            string  `json:"title" validate:"required,min=3,max=255"`
	Slug             string  `json:"slug" validate:"required,min=3,max=100"`
	ShortDescription *string `json:"short_description,omitempty" validate:"omitempty,max=500"`
	Description      string  `json:"description" validate:"required,min=10,max=10000"`
	CoverURL         string  `json:"cover_url" validate:"required,http_url,max=500"`
	Visibility       string  `json:"visibility" validate:"required,oneof=public private"`
	Difficulty       *string `json:"difficulty,omitempty" validate:"omitempty,oneof=beginner intermediate advanced"`
	Language         *string `json:"language,omitempty" validate:"omitempty,min=2,max=10"`
	CategoryID       *int64  `json:"category_id,omitempty" validate:"omitempty,gt=0"`
}

type UpdateCourseRequest struct {
	Title            *string `json:"title,omitempty" validate:"omitempty,min=3,max=255"`
	Slug             *string `json:"slug,omitempty" validate:"omitempty,min=3,max=100"`
	ShortDescription *string `json:"short_description,omitempty" validate:"omitempty,max=500"`
	Description      *string `json:"description,omitempty" validate:"omitempty,min=10,max=10000"`
	CoverURL         *string `json:"cover_url,omitempty" validate:"omitempty,http_url,max=500"`
	Visibility       *string `json:"visibility,omitempty" validate:"omitempty,oneof=public private"`
	Difficulty       *string `json:"difficulty,omitempty" validate:"omitempty,oneof=beginner intermediate advanced"`
	Language         *string `json:"language,omitempty" validate:"omitempty,min=2,max=10"`
	CategoryID       *int64  `json:"category_id,omitempty" validate:"omitempty,gt=0"`
}

type CreateCourseResponseData struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Slug      string    `json:"slug"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateCourseResponse struct {
	Data CreateCourseResponseData `json:"data"`
}

type ListCoursesResponse struct {
	Data PaginatedCourses `json:"data"`
}

type PaginatedCourses struct {
	Courses  []Course `json:"courses"`
	Page     int64    `json:"page"`
	PageSize int64    `json:"page_size"`
	Total    int64    `json:"total"`
}
