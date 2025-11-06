package dto

import (
	"edtech/internal/domain"
	"time"
)

const (
	VisibilityPublic  = "public"
	VisibilityPrivate = "private"
)

const (
	StatusDraft     = "draft"
	StatusPublished = "published"
)

type Course struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	CoverURL    string    `json:"cover_url"`
	CreatedBy   int64     `json:"created_by"`
	Visibility  string    `json:"visibility" validate:"required,oneof=public private"`
	Status      string    `json:"status" validate:"required,oneof=draft published"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
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
	Title       string `json:"title" validate:"required"`
	Slug        string `json:"slug" validate:"required"`
	Description string `json:"description" validate:"required"`
	CoverURL    string `json:"cover_url" validate:"required"`
	CreatedBy   int64  `json:"created_by" validate:"required"`
	Visibility  string `json:"visibility" validate:"required,oneof=public private"`
	Status      string `json:"status" validate:"required,oneof=draft published"`
}

type CreateCourseResponseData struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Slug      string    `json:"slug"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateCourseResponse struct {
	Data    CreateCourseResponseData `json:"data"`
	Message string                   `json:"message"`
}

type ListCoursesResponse struct {
	Data PaginatedCourses `json:"data"`
}

type PaginatedCourses struct {
	Courses  []domain.Course `json:"courses"`
	Page     int64           `json:"page"`
	PageSize int64           `json:"page_size"`
	Total    int64           `json:"total"`
}
