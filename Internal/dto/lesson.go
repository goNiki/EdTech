package dto

import (
	"time"

	"edtech/internal/domain"
)

type Lesson struct {
	ID           int64                `json:"id"`
	CourseID     int64                `json:"course_id"`
	SectionID    *int64               `json:"section_id,omitempty"`
	Title        string               `json:"title"`
	Description  string               `json:"description"`
	CoverURL     string               `json:"cover_url"`
	Content      string               `json:"content"`
	Type         string               `json:"type"` // lecture, exercise, quiz, assignment
	Position     int64                `json:"position"`
	Duration     *int                 `json:"duration,omitempty"` // seconds
	IsFree       bool                 `json:"is_free"`
	Status       string               `json:"status"`
	QuizSettings *domain.QuizSettings `json:"quiz_settings,omitempty"`
	CreatedAt    time.Time            `json:"created_at"`
	UpdatedAt    time.Time            `json:"updated_at"`
	PublishedAt  *time.Time           `json:"published_at,omitempty"`
}

type CreateLessonRequest struct {
	CourseID     int64                `json:"course_id" validate:"required,gt=0"`
	SectionID    *int64               `json:"section_id,omitempty" validate:"omitempty,gt=0"`
	Title        string               `json:"title" validate:"required,min=1,max=255"`
	Description  string               `json:"description,omitempty" validate:"omitempty,max=2000"`
	CoverURL     string               `json:"cover_url,omitempty" validate:"omitempty,max=500"`
	Content      string               `json:"content"`
	Type         string               `json:"type" validate:"required,oneof=lecture exercise practice quiz assignment homework"`
	Duration     *int                 `json:"duration,omitempty" validate:"omitempty,gte=0"`
	IsFree       bool                 `json:"is_free"`
	Status       string               `json:"status,omitempty" validate:"omitempty,oneof=draft published archived"`
	QuizSettings *domain.QuizSettings `json:"quiz_settings,omitempty"`
}

type UpdateLessonRequest struct {
	SectionID    *int64               `json:"section_id,omitempty" validate:"omitempty,gt=0"`
	Title        *string              `json:"title,omitempty" validate:"omitempty,min=1,max=255"`
	Description  *string              `json:"description,omitempty" validate:"omitempty,max=2000"`
	CoverURL     *string              `json:"cover_url,omitempty" validate:"omitempty,max=500"`
	Content      *string              `json:"content,omitempty"`
	Type         *string              `json:"type,omitempty" validate:"omitempty,oneof=lecture exercise practice quiz assignment homework"`
	Duration     *int                 `json:"duration,omitempty" validate:"omitempty,gte=0"`
	IsFree       *bool                `json:"is_free,omitempty"`
	Position     *int64               `json:"position,omitempty"`
	Status       *string              `json:"status,omitempty" validate:"omitempty,oneof=draft published archived"`
	QuizSettings *domain.QuizSettings `json:"quiz_settings,omitempty"`
}

type CreateLessonResponse struct {
	Lesson Lesson `json:"lesson"`
}

type LessonNavigationResponse struct {
	Code    int                     `json:"code"`
	Message string                  `json:"message"`
	Data    LessonNavigationDataDTO `json:"data"`
}

type LessonNavigationDataDTO struct {
	CurrentLesson LessonNavCurrentDTO   `json:"current_lesson"`
	Course        LessonNavCourseDTO    `json:"course"`
	PrevLesson    *LessonNavNeighborDTO `json:"prev_lesson"`
	NextLesson    *LessonNavNeighborDTO `json:"next_lesson"`
	Syllabus      []LessonNavSectionDTO `json:"syllabus"`
}

type LessonNavCurrentDTO struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Position  int64  `json:"position"`
	SectionID *int64 `json:"section_id"`
}

type LessonNavCourseDTO struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	Slug  string `json:"slug"`
}

type LessonNavNeighborDTO struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

type LessonNavSectionDTO struct {
	SectionID    int64              `json:"section_id"`
	SectionTitle string             `json:"section_title"`
	Position     int                `json:"position"`
	Lessons      []LessonNavItemDTO `json:"lessons"`
}

type LessonNavItemDTO struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Position    int64  `json:"position"`
	IsCompleted bool   `json:"is_completed"`
	Score       int    `json:"score"`
}

