package domain

import (
	"regexp"
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
	Id                int64
	Title             string
	Slug              string
	ShortDescription  *string
	Description       string
	CoverURL          string
	IntroVideoURL     *string
	CreatedBy         int64
	Visibility        string
	Status            string
	Difficulty        *string
	Language          *string
	EstimatedDuration *int
	CategoryID        *int64
	TotalLessons      int
	TotalSections     int
	EnrolledCount     int
	CreatedAt         time.Time
	UpdatedAt         time.Time
	PublishedAt       *time.Time
	ArchivedAt        *time.Time
	DeletedAt         *time.Time
}

func (c *Course) Publish() error {
	if c.Status == StatusPublished {
		return errorsAPP.ErrCourseAlredyPublished
	}

	c.Status = StatusPublished

	c.UpdatedAt = time.Now()

	return nil
}

type PaginatedCourses struct {
	Courses  []Course `json:"courses"`
	Page     int64    `json:"page"`
	PageSize int64    `json:"pagesize"`
	Total    int64    `json:"total"`
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

func (c *Course) Validate() error {
	if c.Title == "" {
		return errorsAPP.ErrEmptyTitle
	}
	
	if len(c.Title) > 200 {
		return errorsAPP.ErrTitleTooLong
	}
	
	if c.Slug == "" {
		return errorsAPP.ErrEmptySlug
	}
	
	// Slug validation
	slugRegex := regexp.MustCompile("^[a-z0-9-]+$")
	if !slugRegex.MatchString(c.Slug) {
		return errorsAPP.ErrInvalidSlug
	}
	
	return nil
}