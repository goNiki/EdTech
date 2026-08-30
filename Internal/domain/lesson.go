package domain

import "time"

type Lesson struct {
	ID          int64
	CourseID    int64
	SectionID   *int64
	Title       string
	Description string
	CoverURL    string
	Content     string
	Type        string
	Position    int64
	Duration    *int
	IsFree      bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	PublishedAt *time.Time
	DeletedAt   *time.Time
}
