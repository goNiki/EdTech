package domain

import (
	errorsAPP "edtech/pkg/errors"
	"time"
)

const (
	VisibilityPublic  = "public"
	VisibilityPrivate = "private"
)

const (
	StatusDraft     = "draft"
	StatusPublished = "published"
)

type RoleCourses string

const (
	StudentRole RoleCourses = "student"
	TeacherRole RoleCourses = "teacher"
	CreatorRole RoleCourses = "creator"
)

type Course struct {
	Id          int64     `db:"id"`
	Title       string    `db:"title"`
	Slug        string    `db:"slug"`
	Description string    `db:"description"`
	CoverURL    string    `db:"cover_url"`
	CreatedBy   int64     `db:"created_by"`
	Visibility  string    `db:"visibility"`
	Status      string    `db:"status"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func (c *Course) Publish() error {
	if c.Status == "published" {
		return errorsAPP.ErrCourseAlredyPublished
	}

	c.Status = "published"

	c.UpdatedAt = time.Now()

	return nil
}

type PaginatedCourses struct {
	Courses  []Course `json:"courses"`
	Page     int      `json:"page"`
	PageSize int      `json:"pagesize"`
	Total    int      `json:"total"`
}

type CourseWithLessons struct {
	Course  Course   `json:"course"`
	Lessons []Lesson `json:"lessons"`
}

type EnrolledInCourse struct {
	UserID   int64  `json:"userid"`
	CourseID int64  `json:"courseid"`
	Role     string `json:"role"`
}
