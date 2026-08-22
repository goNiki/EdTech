package domain

import "time"

type Lesson struct {
	ID          int64     `db:"id"`
	CourseID    int64     `db:"course_id"`
	SectionID   *int64    `db:"section_id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	CoverURL    string    `db:"cover_url"`
	Content     string    `db:"content"`
	Type        string    `db:"type"`
	Position    int64     `db:"position"`
	Duration    *int      `db:"duration"`
	IsFree      bool      `db:"is_free"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	PublishedAt *time.Time `db:"published_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
}