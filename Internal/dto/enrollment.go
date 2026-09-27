package dto

import "time"

type TeacherEnrollRequest struct {
	UserID *int64 `json:"user_id,omitempty"`
	Email  string `json:"email,omitempty"`
	Role   string `json:"role,omitempty"`
}

type CourseStudentDTO struct {
	UserID             int64     `json:"user_id"`
	FirstName          string    `json:"first_name"`
	LastName           string    `json:"last_name"`
	Username           string    `json:"username"`
	Email              string    `json:"email"`
	AvatarURL          *string   `json:"avatar_url,omitempty"`
	Role               string    `json:"role"`
	EnrolledAt         time.Time `json:"enrolled_at"`
	ProgressPercentage float64   `json:"progress_percentage"`
	CompletedLessons   int       `json:"completed_lessons"`
	TotalLessons       int       `json:"total_lessons"`
	AverageScore       float64   `json:"average_score"`
	HasPending         bool      `json:"has_pending_homeworks"`
}

type PaginatedStudentsResponse struct {
	Students []CourseStudentDTO `json:"students"`
	Total    int64              `json:"total"`
	Page     int64              `json:"page"`
	PageSize int64              `json:"page_size"`
}
