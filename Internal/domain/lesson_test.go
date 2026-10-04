package domain_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"edtech/internal/domain"
	errorsAPP "edtech/pkg/errors"
)

func TestValidateLessonContent_EmptyOrSpaces(t *testing.T) {
	if err := domain.ValidateLessonContent(""); err != nil {
		t.Fatalf("expected nil for empty content, got %v", err)
	}
	if err := domain.ValidateLessonContent("   \n\t "); err != nil {
		t.Fatalf("expected nil for whitespace content, got %v", err)
	}
}

func TestValidateLessonContent_InvalidJSON(t *testing.T) {
	err := domain.ValidateLessonContent("{ invalid json content }")
	if err == nil {
		t.Fatal("expected error for malformed json, got nil")
	}
	if !errors.Is(err, errorsAPP.ErrLessonValidation) {
		t.Fatalf("expected ErrLessonValidation, got %v", err)
	}
}

func TestValidateLessonContent_InvalidRoot(t *testing.T) {
	err := domain.ValidateLessonContent(`"just a string"`)
	if err == nil {
		t.Fatal("expected error for non-object/array root, got nil")
	}
	if !errors.Is(err, errorsAPP.ErrLessonValidation) {
		t.Fatalf("expected ErrLessonValidation, got %v", err)
	}
}

func TestValidateLessonContent_ValidPuckFormat_With50Quizzes(t *testing.T) {
	blocks := make([]map[string]interface{}, 60)
	for i := 0; i < 60; i++ {
		blocks[i] = map[string]interface{}{
			"type": "QuizSingleBlock",
			"props": map[string]interface{}{
				"id":       fmt.Sprintf("quiz-block-%d", i),
				"question": fmt.Sprintf("Вопрос №%d?", i),
				"points":   5,
			},
		}
	}

	puckData := map[string]interface{}{
		"content": blocks,
		"root":    map[string]interface{}{"title": "Урок по Go"},
	}

	bytesData, err := json.Marshal(puckData)
	if err != nil {
		t.Fatalf("failed to marshal puck data: %v", err)
	}

	lesson := &domain.Lesson{
		CourseID:    1,
		Title:       "Введение в Concurrency",
		Description: "Описание урока",
		Content:     string(bytesData),
	}

	if err := lesson.Validate(); err != nil {
		t.Fatalf("expected lesson with 60 quizzes to be valid, got: %v", err)
	}
}

func TestValidateLessonContent_DuplicateBlockIDs(t *testing.T) {
	blocks := []map[string]interface{}{
		{
			"type": "QuizSingleBlock",
			"props": map[string]interface{}{
				"id":       "duplicate-id-1",
				"question": "Вопрос 1",
			},
		},
		{
			"type": "QuizMultiBlock",
			"props": map[string]interface{}{
				"id":       "duplicate-id-1", // duplicate!
				"question": "Вопрос 2",
			},
		},
	}

	puckData := map[string]interface{}{
		"content": blocks,
	}

	bytesData, _ := json.Marshal(puckData)

	lesson := &domain.Lesson{
		CourseID:    1,
		Title:       "Тест",
		Description: "Описание",
		Content:     string(bytesData),
	}

	err := lesson.Validate()
	if err == nil {
		t.Fatal("expected error for duplicate block IDs, got nil")
	}
	if !errors.Is(err, errorsAPP.ErrLessonValidation) {
		t.Fatalf("expected ErrLessonValidation, got %v", err)
	}
}

func TestValidateLessonContent_MissingBlockType(t *testing.T) {
	blocks := []map[string]interface{}{
		{
			"props": map[string]interface{}{
				"id": "quiz-1",
			},
		},
	}

	puckData := map[string]interface{}{
		"content": blocks,
	}
	bytesData, _ := json.Marshal(puckData)

	err := domain.ValidateLessonContent(string(bytesData))
	if err == nil {
		t.Fatal("expected error for missing type, got nil")
	}
	if !errors.Is(err, errorsAPP.ErrLessonValidation) {
		t.Fatalf("expected ErrLessonValidation, got %v", err)
	}
}

func TestValidateLessonContent_QuizMissingPropsID(t *testing.T) {
	blocks := []map[string]interface{}{
		{
			"type": "QuizSingleBlock",
			"props": map[string]interface{}{
				"question": "Вопрос без ID",
			},
		},
	}

	puckData := map[string]interface{}{
		"content": blocks,
	}
	bytesData, _ := json.Marshal(puckData)

	err := domain.ValidateLessonContent(string(bytesData))
	if err == nil {
		t.Fatal("expected error for quiz block missing ID, got nil")
	}
	if !errors.Is(err, errorsAPP.ErrLessonValidation) {
		t.Fatalf("expected ErrLessonValidation, got %v", err)
	}
}
