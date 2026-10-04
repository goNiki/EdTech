package lesson_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"edtech/internal/domain"
	lessonHandler "edtech/internal/interfaces/handlers/lesson"
	"edtech/internal/interfaces/middleware/auth"
	"edtech/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type mockLessonService struct {
	service.LessonServices
	updatedLesson *domain.Lesson
	updateErr     error
}

func (m *mockLessonService) UpdateLesson(ctx context.Context, userID int64, lesson *domain.Lesson) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	m.updatedLesson = lesson
	return nil
}

type mockAuthMiddleware struct {
	auth.AuthMiddleware
	userID int64
}

func (m *mockAuthMiddleware) GetUserID(ctx context.Context) int64 {
	return m.userID
}

func TestUpdateLesson_ValidContent_With50Quizzes(t *testing.T) {
	mockSvc := &mockLessonService{}
	authMw := &mockAuthMiddleware{userID: 1}
	val := validator.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	h := lessonHandler.NewLessonHandler(mockSvc, logger, val, authMw)

	// Сгенерируем 55 валидных квиз-блоков
	blocks := make([]map[string]interface{}, 55)
	for i := 0; i < 55; i++ {
		blocks[i] = map[string]interface{}{
			"type": "QuizSingleBlock",
			"props": map[string]interface{}{
				"id":       fmt.Sprintf("quiz-%d", i),
				"question": fmt.Sprintf("Вопрос №%d?", i),
			},
		}
	}
	puckJSON, _ := json.Marshal(map[string]interface{}{"content": blocks})
	contentStr := string(puckJSON)

	reqBody := map[string]interface{}{
		"title":   "Урок с тестами",
		"content": contentStr,
	}
	bodyBytes, _ := json.Marshal(reqBody)

	r := chi.NewRouter()
	r.Patch("/api/v1/lessons/{id}", h.UpdateLesson)

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/lessons/42", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if mockSvc.updatedLesson == nil || mockSvc.updatedLesson.Content != contentStr {
		t.Fatalf("expected lesson content to be updated")
	}
}

func TestUpdateLesson_MalformedJSON_Rejected(t *testing.T) {
	mockSvc := &mockLessonService{}
	authMw := &mockAuthMiddleware{userID: 1}
	val := validator.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	h := lessonHandler.NewLessonHandler(mockSvc, logger, val, authMw)

	reqBody := map[string]interface{}{
		"title":   "Урок с битым JSON",
		"content": "{ broken: json, ",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	r := chi.NewRouter()
	r.Patch("/api/v1/lessons/{id}", h.UpdateLesson)

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/lessons/42", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d (body: %s)", rec.Code, rec.Body.String())
	}
}

func TestUpdateLesson_MaxBytesReader_Protection(t *testing.T) {
	mockSvc := &mockLessonService{}
	authMw := &mockAuthMiddleware{userID: 1}
	val := validator.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	h := lessonHandler.NewLessonHandler(mockSvc, logger, val, authMw)

	// Создаем тело > 5 МБ
	hugeContent := strings.Repeat("A", 6*1024*1024)
	reqBody := fmt.Sprintf(`{"title":"Huge","content":"%s"}`, hugeContent)

	r := chi.NewRouter()
	r.Patch("/api/v1/lessons/{id}", h.UpdateLesson)

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/lessons/42", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request for >5MB payload, got %d", rec.Code)
	}
}

func TestUpdateLesson_Unauthorized(t *testing.T) {
	mockSvc := &mockLessonService{}
	authMw := &mockAuthMiddleware{userID: 0}
	val := validator.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	h := lessonHandler.NewLessonHandler(mockSvc, logger, val, authMw)

	r := chi.NewRouter()
	r.Patch("/api/v1/lessons/{id}", h.UpdateLesson)

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/lessons/42", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized, got %d", rec.Code)
	}
}
