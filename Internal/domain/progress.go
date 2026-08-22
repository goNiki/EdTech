package domain

import "time"

type ProgressStatus string

const (
	ProgressStatusNotStarted ProgressStatus = "not_started"
	ProgressStatusInProgress ProgressStatus = "in_progress"
	ProgressStatusCompleted  ProgressStatus = "completed"
)

type LessonProgress struct {
	ID          int64          `db:"id"`
	UserID      int64          `db:"user_id"`
	LessonID    int64          `db:"lesson_id"`
	Status      ProgressStatus `db:"status"`
	Score       *int           `db:"score"`
	TimeSpent   int            `db:"time_spent"`
	LastPos     int            `db:"last_position"`
	StartedAt   *time.Time     `db:"started_at"`
	CompletedAt *time.Time     `db:"completed_at"`
	UpdatedAt   time.Time      `db:"updated_at"`
}

type CourseProgress struct {
	ID             int64          `db:"id"`
	UserID         int64          `db:"user_id"`
	CourseID       int64          `db:"course_id"`
	Status         ProgressStatus `db:"status"`
	Percent        int            `db:"completion_percentage"`
	CompletedLess  int            `db:"completed_lessons"`
	StartedAt      *time.Time     `db:"started_at"`
	CompletedAt    *time.Time     `db:"completed_at"`
	LastAccessedAt time.Time      `db:"last_accessed_at"`
}
