package domain

import "time"

type Section struct {
	ID          int64      `json:"id"`
	CourseID    int64      `json:"course_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Position    int        `json:"position"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type SectionWithLessons struct {
	Section Section  `json:"section"`
	Lessons []Lesson `json:"lessons"`
}

type CourseStructure struct {
	Course             Course               `json:"course"`
	Sections           []SectionWithLessons `json:"sections"`
	UnsectionedLessons []Lesson             `json:"unsectioned_lessons,omitempty"`
}
