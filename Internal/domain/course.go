package domain

import (
	errorsAPP "edtech/pkg/errors"
	"regexp"
	"time"
)

const (
	VisibilityPublic  = "public"
	VisibilityPrivate = "private"
)

const (
	StatusDraft     = "draft"
	StatusPublished = "published"
	StatusArchived   = "archived"
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
	Rating            float64
	ReviewsCount      int
	CreatedAt         time.Time
	UpdatedAt         time.Time
	PublishedAt       *time.Time
	ArchivedAt        *time.Time
	DeletedAt         *time.Time
}

func (c *Course) CanArchive() error {
	if c.Status == StatusArchived {
		return errorsAPP.ErrCourseAlreadyArchived
	}
	return nil
}

func (c *Course) Archive(time time.Time) {
	c.Status = StatusArchived
	c.UpdatedAt = time
	c.ArchivedAt = &time
}

func (c *Course) CanPublish() error {
	if c.Status == StatusPublished {
		return errorsAPP.ErrCourseAlreadyPublished
	}
	return nil
}

func (c *Course) Publish(time time.Time) {

	c.Status = StatusPublished

	c.UpdatedAt = time

	c.PublishedAt = &time

}

func (c *Course) CanSelfEnroll() error {
	if c.Status != StatusPublished || c.Visibility != VisibilityPublic {
		return errorsAPP.ErrForbidden
	}
	return nil
}

type PaginatedCourses struct {
	Courses  []Course `json:"courses"`
	Page     int64    `json:"page"`
	PageSize int64    `json:"pagesize"`
	Total    int64    `json:"total"`
}

type CourseFilter struct {
	Search     *string
	CategoryID *int64
	CreatedBy  *int64
	Language   *string
	Difficulty *string
	SortBy     string
	SortOrder  string
}

type InputListMyCourse struct {
	UserID     int64
	Role       string
	Pagination Pagination
	Filter     CourseFilter
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

type UpdateCourseInput struct {
	Title             *string
	Slug              *string
	ShortDescription  *string
	Description       *string
	CoverURL          *string
	IntroVideoURL     *string
	Visibility        *string
	Difficulty        *string
	Language          *string
	EstimatedDuration *int
	CategoryID        *int64
}

func (c *Course) Update(input UpdateCourseInput) {
	if input.Title != nil {
		c.Title = *input.Title
	}
	if input.Slug != nil {
		c.Slug = *input.Slug
	}
	if input.ShortDescription != nil {
		c.ShortDescription = input.ShortDescription
	}
	if input.Description != nil {
		c.Description = *input.Description
	}
	if input.CoverURL != nil {
		c.CoverURL = *input.CoverURL
	}
	if input.IntroVideoURL != nil {
		c.IntroVideoURL = input.IntroVideoURL
	}
	if input.Visibility != nil {
		c.Visibility = *input.Visibility
	}
	if input.Difficulty != nil {
		c.Difficulty = input.Difficulty
	}
	if input.Language != nil {
		c.Language = input.Language
	}
	if input.EstimatedDuration != nil {
		c.EstimatedDuration = input.EstimatedDuration
	}
	if input.CategoryID != nil {
		c.CategoryID = input.CategoryID
	}
}
