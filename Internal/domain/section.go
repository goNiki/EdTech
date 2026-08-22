package domain

import "time"

type Section struct {
	ID          int64      `db:"id"`
	CourseID    int64      `db:"course_id"`
	Title       string     `db:"title"`
	Description string     `db:"description"`
	Position    int        `db:"position"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
}
