package domain

import (
	errorsAPP "edtech/pkg/errors"
	"time"
)

type Course struct {
	Id          int       `db:"id"`
	Title       string    `db:"title"`
	Slug        string    `db:"slug"`
	Description string    `db:"description"`
	CoverURL    string    `db:"cover_url"`
	CreatedBy   int       `db:"created_by"`
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
	UserID   int    `json:"userid"`
	CourseID int    `json:"courseid"`
	Role     string `json:"role"`
}
