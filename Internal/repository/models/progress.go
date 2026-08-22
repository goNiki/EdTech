package models

import "time"

type LessonProgress struct {
	ID          int64
	UserID      int64
	LessonID    int64
	Status      string
	Score       *int
	TimeSpent   int
	LastPos     int
	StartedAt   *time.Time
	CompletedAt *time.Time
	UpdatedAt   time.Time
}

type CourseProgress struct {
	ID             int64
	UserID         int64
	CourseID       int64
	Status         string
	Percent        int
	CompletedLess  int
	StartedAt      *time.Time
	CompletedAt    *time.Time
	LastAccessedAt time.Time
}
