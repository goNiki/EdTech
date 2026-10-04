package domain

import "time"

type Review struct {
	ID        int64     `json:"id"`
	CourseID  int64     `json:"course_id"`
	UserID    int64     `json:"user_id"`
	Rating    int       `json:"rating"`
	Comment   *string   `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Enriched user details (for display in reviews list)
	UserName   string `json:"user_name,omitempty"`
	UserAvatar string `json:"user_avatar,omitempty"`
}

type CourseReviewsSummary struct {
	Reviews       []Review `json:"reviews"`
	AverageRating float64  `json:"average_rating"`
	ReviewsCount  int      `json:"reviews_count"`
	Total         int64    `json:"total"`
}
