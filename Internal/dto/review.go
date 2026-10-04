package dto

import "time"

type CreateReviewRequest struct {
	Rating  int     `json:"rating" validate:"required,min=1,max=5"`
	Comment *string `json:"comment,omitempty"`
}

type ReviewResponse struct {
	ID         int64     `json:"id"`
	CourseID   int64     `json:"course_id"`
	UserID     int64     `json:"user_id"`
	Rating     int       `json:"rating"`
	Comment    *string   `json:"comment,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	UserName   string    `json:"user_name,omitempty"`
	UserAvatar string    `json:"user_avatar,omitempty"`
}

type CourseReviewsListResponse struct {
	Reviews       []ReviewResponse `json:"reviews"`
	AverageRating float64          `json:"average_rating"`
	ReviewsCount  int              `json:"reviews_count"`
	Total         int64            `json:"total"`
}
