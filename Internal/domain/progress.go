package domain

import "time"

type ProgressStatus string

const (
	ProgressStatusNotStarted ProgressStatus = "not_started"
	ProgressStatusInProgress ProgressStatus = "in_progress"
	ProgressStatusCompleted  ProgressStatus = "completed"
)

type LessonProgress struct {
	ID          int64
	UserID      int64
	LessonID    int64
	CourseID    int64
	Status      ProgressStatus
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
	CompletedLess  int
	TotalLessons   int
	Percent        int
	TotalWatchTime int
	AverageScore   *float64
	StartedAt      *time.Time
	LastAccessedAt time.Time
	CompletedAt    *time.Time
}

type UpdateProgressInput struct {
	Status       ProgressStatus `json:"status"`
	TimeSpent    int            `json:"time_spent"`
	LastPosition int            `json:"last_position"`
}
