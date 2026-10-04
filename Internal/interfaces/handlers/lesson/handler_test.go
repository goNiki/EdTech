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
	errorsAPP "edtech/pkg/errors"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type mockLessonService struct {
	service.LessonServices
	updatedLesson *domain.Lesson
	updateErr     error
	navContext    *domain.LessonNavigationContext
	navErr        error
}

func (m *mockLessonService) UpdateLesson(ctx context.Context, userID int64, lesson *domain.Lesson) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	m.updatedLesson = lesson
	return nil
}

func (m *mockLessonService) GetLessonNavigationContext(ctx context.Context, userID int64, lessonID int64) (*domain.LessonNavigationContext, error) {
	if m.navErr != nil {
		return nil, m.navErr
	}
	if m.navContext != nil {
		return m.navContext, nil
	}
	return &domain.LessonNavigationContext{
		CurrentLesson: domain.LessonNavCurrent{ID: lessonID, Title: "Урок", Position: 1},
		Course:        domain.LessonNavCourse{ID: 10, Title: "Курс", Slug: "kurs"},
	}, nil
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

func TestGetLessonNavigationContext_Success(t *testing.T) {
	prev := &domain.LessonNavNeighbor{ID: 41, Title: "Предыдущий"}
	next := &domain.LessonNavNeighbor{ID: 43, Title: "Следующий"}
	secID := int64(12)
	navCtx := &domain.LessonNavigationContext{
		CurrentLesson: domain.LessonNavCurrent{ID: 42, Title: "Текущий", Position: 2, SectionID: &secID},
		Course:        domain.LessonNavCourse{ID: 10, Title: "Курс", Slug: "kurs"},
		PrevLesson:    prev,
		NextLesson:    next,
		Syllabus: []domain.LessonNavSection{
			{
				SectionID:    12,
				SectionTitle: "Модуль 1",
				Position:     1,
				Lessons: []domain.LessonNavItem{
					{ID: 41, Title: "Предыдущий", Position: 1, IsCompleted: true, Score: 100},
					{ID: 42, Title: "Текущий", Position: 2, IsCompleted: false, Score: 0},
					{ID: 43, Title: "Следующий", Position: 3, IsCompleted: false, Score: 0},
				},
			},
		},
	}

	mockSvc := &mockLessonService{navContext: navCtx}
	authMw := &mockAuthMiddleware{userID: 1}
	val := validator.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	h := lessonHandler.NewLessonHandler(mockSvc, logger, val, authMw)

	r := chi.NewRouter()
	r.Get("/api/v1/lessons/{id}/navigation", h.GetLessonNavigationContext)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lessons/42/navigation", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d, body: %s", rec.Code, rec.Body.String())
	}

	var resp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			CurrentLesson struct {
				ID int64 `json:"id"`
			} `json:"current_lesson"`
			PrevLesson *struct {
				ID int64 `json:"id"`
			} `json:"prev_lesson"`
			NextLesson *struct {
				ID int64 `json:"id"`
			} `json:"next_lesson"`
			Syllabus []struct {
				SectionID int64 `json:"section_id"`
			} `json:"syllabus"`
		} `json:"data"`
	}

	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Code != 200 {
		t.Errorf("expected code 200, got %d", resp.Code)
	}
	if resp.Data.CurrentLesson.ID != 42 {
		t.Errorf("expected current lesson 42, got %d", resp.Data.CurrentLesson.ID)
	}
	if resp.Data.PrevLesson == nil || resp.Data.PrevLesson.ID != 41 {
		t.Errorf("expected prev lesson 41, got %+v", resp.Data.PrevLesson)
	}
	if resp.Data.NextLesson == nil || resp.Data.NextLesson.ID != 43 {
		t.Errorf("expected next lesson 43, got %+v", resp.Data.NextLesson)
	}
	if len(resp.Data.Syllabus) != 1 || resp.Data.Syllabus[0].SectionID != 12 {
		t.Errorf("expected 1 syllabus section with id 12, got %+v", resp.Data.Syllabus)
	}
}

func TestGetLessonNavigationContext_Unauthorized(t *testing.T) {
	mockSvc := &mockLessonService{}
	authMw := &mockAuthMiddleware{userID: 0}
	val := validator.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	h := lessonHandler.NewLessonHandler(mockSvc, logger, val, authMw)

	r := chi.NewRouter()
	r.Get("/api/v1/lessons/{id}/navigation", h.GetLessonNavigationContext)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lessons/42/navigation", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized, got %d", rec.Code)
	}
}

func TestGetLessonNavigationContext_Forbidden(t *testing.T) {
	mockSvc := &mockLessonService{navErr: errorsAPP.ErrForbidden}
	authMw := &mockAuthMiddleware{userID: 99}
	val := validator.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	h := lessonHandler.NewLessonHandler(mockSvc, logger, val, authMw)

	r := chi.NewRouter()
	r.Get("/api/v1/lessons/{id}/navigation", h.GetLessonNavigationContext)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lessons/42/navigation", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden, got %d", rec.Code)
	}
}

func TestGetLessonNavigationContext_InvalidID(t *testing.T) {
	mockSvc := &mockLessonService{}
	authMw := &mockAuthMiddleware{userID: 1}
	val := validator.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	h := lessonHandler.NewLessonHandler(mockSvc, logger, val, authMw)

	r := chi.NewRouter()
	r.Get("/api/v1/lessons/{id}/navigation", h.GetLessonNavigationContext)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lessons/not-a-number/navigation", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d", rec.Code)
	}
}

