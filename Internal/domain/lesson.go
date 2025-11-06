package domain

import "time"

type Lesson struct {
	ID          int64     `db:"id"`
	CourseID    int64     `db:"course_id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	CoverURL    string    `db:"cover_url"`
	Content     string    `db:"content"`
	Position    int64     `db:"position"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
