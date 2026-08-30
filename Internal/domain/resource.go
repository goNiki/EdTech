package domain

import "time"

type Resource struct {
	ID            int64
	LessonID      *int64
	CourseID      *int64
	Title         *string
	Description   *string
	Type          string
	Path          string
	Mime          string
	Size          int64
	ExternalURL   *string
	Duration      *int
	OrderPosition int
	IsRequired    bool
	CreatedAt     time.Time
	DeletedAt     *time.Time
}
