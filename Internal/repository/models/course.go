package models

import "time"

type Course struct {
	ID                int64
	Title             string
	Slug              string
	ShortDescription  *string
	Description       string
	CoverURL          string
	IntroVideoURL     *string
	CreatedBy         int64
	Visibility        string
	Status            string
	Difficulty        *string
	Language          *string
	EstimatedDuration *int
	CategoryID        *int64
	TotalLessons      int
	TotalSections     int
	EnrolledCount     int
	CreatedAt         time.Time
	UpdatedAt         time.Time
	PublishedAt       *time.Time
	ArchivedAt        *time.Time
	DeletedAt         *time.Time
}
