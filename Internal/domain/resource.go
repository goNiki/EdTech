package domain

import "time"

type Resource struct {
	ID            int64     `db:"id"`
	LessonID      *int64    `db:"lesson_id"`
	CourseID      *int64    `db:"course_id"`
	Title         *string   `db:"title"`
	Description   *string   `db:"description"`
	Type          string    `db:"type"`
	Path          string    `db:"path"`
	Mime          string    `db:"mime"`
	Size          int64     `db:"size"`
	ExternalURL   *string   `db:"external_url"`
	Duration      *int      `db:"duration"`
	OrderPosition int       `db:"order_position"`
	IsRequired    bool      `db:"is_required"`
	CreatedAt     time.Time `db:"created_at"`
	DeletedAt     *time.Time `db:"deleted_at"`
}