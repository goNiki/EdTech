package domain

import (
	errorsAPP "edtech/pkg/errors"
	"time"
)

type Lesson struct {
	ID          int64
	CourseID    int64
	SectionID   *int64
	Title       string
	Description string
	CoverURL    string
	Content     string
	Type        string
	Position    int64
	Duration    *int
	IsFree      bool
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	PublishedAt *time.Time
	DeletedAt   *time.Time
}

func (l *Lesson) CanPublish() error {
	if l.Status == StatusPublished {
		return errorsAPP.ErrLessonAlreadyPublished
	}
	return nil
}

func (l *Lesson) Publish(now time.Time) {
	l.Status = StatusPublished
	l.PublishedAt = &now
	l.UpdatedAt = now
}

func (l *Lesson) CanArchive() error {
	if l.Status == StatusArchived {
		return errorsAPP.ErrLessonAlreadyArchived
	}
	return nil
}

func (l *Lesson) Archive(now time.Time) {
	l.Status = StatusArchived
	l.UpdatedAt = now
}

func (l *Lesson) IsPublished() bool {
	return l.Status == StatusPublished
}
