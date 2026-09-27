package models

import "time"

type Section struct {
	ID          int64
	CourseID    int64
	Title       string
	Description string
	Position    int
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
