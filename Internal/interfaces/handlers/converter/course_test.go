package converter_test

import (
	"testing"
	"time"

	"edtech/internal/domain"
	"edtech/internal/interfaces/handlers/converter"
)

func TestCourseToDTO_WithAuthor(t *testing.T) {
	avatar := "https://example.com/avatar.jpg"
	now := time.Now()

	course := &domain.Course{
		Id:          10,
		Title:       "Курс подготовки к ЕГЭ",
		Slug:        "ege-prep",
		Description: "Полное описание курса",
		CreatedBy:   5,
		Author: &domain.CourseAuthorInfo{
			ID:            5,
			Name:          "Иван Иванов",
			AvatarURL:     &avatar,
			Headline:      "Старший эксперт ЕГЭ",
			Bio:           "Преподаватель с 10-летним стажем",
			CoursesCount:  4,
			TotalStudents: 1250,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	dto := converter.CourseToDTO(course)

	if dto.ID != 10 {
		t.Errorf("expected ID 10, got %d", dto.ID)
	}
	if dto.Author == nil {
		t.Fatal("expected Author to be non-nil")
	}
	if dto.Author.ID != 5 {
		t.Errorf("expected Author ID 5, got %d", dto.Author.ID)
	}
	if dto.Author.Name != "Иван Иванов" {
		t.Errorf("expected Author Name 'Иван Иванов', got %q", dto.Author.Name)
	}
	if dto.Author.Headline != "Старший эксперт ЕГЭ" {
		t.Errorf("expected Author Headline 'Старший эксперт ЕГЭ', got %q", dto.Author.Headline)
	}
	if dto.Author.Bio != "Преподаватель с 10-летним стажем" {
		t.Errorf("expected Author Bio 'Преподаватель с 10-летним стажем', got %q", dto.Author.Bio)
	}
	if dto.Author.CoursesCount != 4 {
		t.Errorf("expected Author CoursesCount 4, got %d", dto.Author.CoursesCount)
	}
	if dto.Author.TotalStudents != 1250 {
		t.Errorf("expected Author TotalStudents 1250, got %d", dto.Author.TotalStudents)
	}
	if dto.Author.AvatarURL == nil || *dto.Author.AvatarURL != avatar {
		t.Errorf("expected Author AvatarURL %q, got %v", avatar, dto.Author.AvatarURL)
	}
}

func TestCourseToDTO_NilAuthor(t *testing.T) {
	course := &domain.Course{
		Id:        11,
		Title:     "Курс без автора",
		Slug:      "no-author",
		Author:    nil,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	dto := converter.CourseToDTO(course)
	if dto.Author != nil {
		t.Errorf("expected Author to be nil, got %+v", dto.Author)
	}
}
