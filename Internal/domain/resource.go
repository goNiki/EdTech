package domain

import "time"

type Resource struct {
	ID        int64     `db:"id"`
	LessonID  int64     `db:"lesson_id"`
	Type      string    `db:"type"`
	Path      string    `db:"path"`
	Mime      string    `db:"mime"`
	Size      int64     `db:"size"`
	CreatedAt time.Time `db:"created_at"`
}
