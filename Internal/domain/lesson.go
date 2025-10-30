package domain

import "time"

type Lesson struct {
	ID          int       `db:"id"`
	CourseID    int       `db:"course_id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	CoverURL    string    `db:"cover_url"`
	Content     string    `db:"content"`
	Position    int       `db:"position"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
