package domain

import "time"

type Resource struct {
	ID        int       `db:"id"`
	LessonID  int       `db:"lesson_id"`
	Type      string    `db:"type"`
	Path      string    `db:"path"`
	Mime      string    `db:"mime"`
	Size      int       `db:"size"`
	CreatedAt time.Time `db:"created_at"`
}
