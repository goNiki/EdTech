package domain

import (
	errorsAPP "edtech/pkg/errors"
	"time"
)

type Section struct {
	ID          int64
	CourseID    int64
	Title       string
	Description string
	Position    int
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

func (s *Section) CanPublish() error {
	if s.Status == StatusPublished {
		return errorsAPP.ErrSectionAlreadyPublished
	}
	return nil
}

func (s *Section) Publish(now time.Time) {
	s.Status = StatusPublished
	s.UpdatedAt = now
}

func (s *Section) CanArchive() error {
	if s.Status == StatusArchived {
		return errorsAPP.ErrSectionAlreadyArchived
	}
	return nil
}

func (s *Section) Archive(now time.Time) {
	s.Status = StatusArchived
	s.UpdatedAt = now
}

func (s *Section) IsPublished() bool {
	return s.Status == StatusPublished
}

type SectionWithLessons struct {
	Section Section
	Lessons []Lesson
}

type CourseStructure struct {
	Course             Course
	Sections           []SectionWithLessons
	UnsectionedLessons []Lesson
}
