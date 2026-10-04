package analytics_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"edtech/internal/domain"
	"edtech/internal/dto"
	analyticsHandler "edtech/internal/interfaces/handlers/analytics"
	"edtech/internal/interfaces/middleware/auth"
	"edtech/internal/service"
	errorsAPP "edtech/pkg/errors"

	"github.com/go-chi/chi/v5"
)

type mockAnalyticsService struct {
	service.AnalyticsServices
	result         *domain.TeacherPendingHomeworksResult
	exportData     []byte
	exportFilename string
	err            error
	calledCID      int64
}

func (m *mockAnalyticsService) ListTeacherPendingHomeworks(ctx context.Context, teacherID, courseID, page, pageSize int64) (*domain.TeacherPendingHomeworksResult, error) {
	m.calledCID = courseID
	if m.err != nil {
		return nil, m.err
	}
	return m.result, nil
}

func (m *mockAnalyticsService) ExportCourseGradebookCSV(ctx context.Context, teacherID, courseID int64) ([]byte, string, error) {
	m.calledCID = courseID
	if m.err != nil {
		return nil, "", m.err
	}
	return m.exportData, m.exportFilename, nil
}

type mockAuthMiddleware struct {
	auth.AuthMiddleware
	userID   int64
	userRole string
}

func (m *mockAuthMiddleware) GetUserID(ctx context.Context) int64 {
	return m.userID
}

func (m *mockAuthMiddleware) GetUserRole(ctx context.Context) string {
	return m.userRole
}

func TestListTeacherPendingHomeworks_Unauthorized(t *testing.T) {
	svc := &mockAnalyticsService{}
	mw := &mockAuthMiddleware{userID: 0}
	h := analyticsHandler.NewAnalyticsHandler(svc, mw)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/teacher/grading/pending", nil)
	rec := httptest.NewRecorder()

	h.ListTeacherPendingHomeworks(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized, got %d", rec.Code)
	}
}

func TestListTeacherPendingHomeworks_ForbiddenRole(t *testing.T) {
	svc := &mockAnalyticsService{}
	mw := &mockAuthMiddleware{userID: 1, userRole: "student"}
	h := analyticsHandler.NewAnalyticsHandler(svc, mw)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/teacher/grading/pending", nil)
	rec := httptest.NewRecorder()

	h.ListTeacherPendingHomeworks(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden for student, got %d", rec.Code)
	}
}

func TestListTeacherPendingHomeworks_Success(t *testing.T) {
	sampleTime := time.Date(2026, 10, 4, 10, 15, 0, 0, time.UTC)
	svc := &mockAnalyticsService{
		result: &domain.TeacherPendingHomeworksResult{
			Items: []domain.PendingHomeworkItem{
				{
					AttemptID:       14,
					AnswerID:        42,
					StudentID:       5,
					StudentName:     "Иван Смирнов",
					StudentEmail:    "student@example.com",
					StudentUsername: "smirnov",
					CourseID:        2,
					CourseTitle:     "Русский язык",
					LessonID:        10,
					LessonTitle:     "Синтаксис и пунктуация",
					QuestionID:      7,
					QuestionText:    "Напишите развернутое сочинение-рассуждение",
					StudentAnswer:   "Текст работы студента...",
					MaxPoints:       25,
					SubmittedAt:     sampleTime,
				},
			},
			Total: 1,
			CoursesSummary: []domain.CoursePendingSummaryItem{
				{CourseID: 2, CourseTitle: "Русский язык", PendingCount: 1},
				{CourseID: 5, CourseTitle: "Основы Go", PendingCount: 0},
			},
		},
	}
	mw := &mockAuthMiddleware{userID: 1, userRole: "teacher"}
	h := analyticsHandler.NewAnalyticsHandler(svc, mw)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/teacher/grading/pending?course_id=2&page=1&page_size=20", nil)
	rec := httptest.NewRecorder()

	h.ListTeacherPendingHomeworks(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Code)
	}

	if svc.calledCID != 2 {
		t.Errorf("expected courseID 2, got %d", svc.calledCID)
	}

	var resp dto.PaginatedTeacherPendingHomeworksResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Total != 1 {
		t.Errorf("expected total 1, got %d", resp.Total)
	}
	if len(resp.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(resp.Items))
	}
	if resp.Items[0].StudentName != "Иван Смирнов" {
		t.Errorf("expected student name 'Иван Смирнов', got %q", resp.Items[0].StudentName)
	}
	if len(resp.CoursesSummary) != 2 {
		t.Fatalf("expected 2 summary items, got %d", len(resp.CoursesSummary))
	}
	if resp.CoursesSummary[0].PendingCount != 1 {
		t.Errorf("expected pending count 1, got %d", resp.CoursesSummary[0].PendingCount)
	}
}

func TestListTeacherPendingHomeworks_ForbiddenCourseAccess(t *testing.T) {
	svc := &mockAnalyticsService{
		err: errorsAPP.ErrForbidden,
	}
	mw := &mockAuthMiddleware{userID: 1, userRole: "teacher"}
	h := analyticsHandler.NewAnalyticsHandler(svc, mw)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/teacher/grading/pending?course_id=99", nil)
	rec := httptest.NewRecorder()

	h.ListTeacherPendingHomeworks(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden, got %d", rec.Code)
	}
}

func TestExportGradebook_Success(t *testing.T) {
	csvData := []byte{0xEF, 0xBB, 0xBF}
	csvData = append(csvData, []byte("ID;Студент;Email\n1;Иван;ivan@mail.ru")...)
	filename := "gradebook_course_13_2026-10-04.csv"

	svc := &mockAnalyticsService{
		exportData:     csvData,
		exportFilename: filename,
	}
	mw := &mockAuthMiddleware{userID: 1, userRole: "teacher"}
	h := analyticsHandler.NewAnalyticsHandler(svc, mw)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/courses/13/analytics/export?format=csv", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("courseid", "13")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rec := httptest.NewRecorder()
	h.ExportGradebook(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d, body: %s", rec.Code, rec.Body.String())
	}

	contentType := rec.Header().Get("Content-Type")
	if contentType != "text/csv; charset=utf-8" {
		t.Errorf("expected Content-Type 'text/csv; charset=utf-8', got %q", contentType)
	}

	contentDisposition := rec.Header().Get("Content-Disposition")
	expectedDisp := `attachment; filename="gradebook_course_13_2026-10-04.csv"`
	if contentDisposition != expectedDisp {
		t.Errorf("expected Content-Disposition %q, got %q", expectedDisp, contentDisposition)
	}

	cacheControl := rec.Header().Get("Cache-Control")
	if cacheControl != "no-cache" {
		t.Errorf("expected Cache-Control 'no-cache', got %q", cacheControl)
	}

	if rec.Body.String() != string(csvData) {
		t.Errorf("expected body %q, got %q", string(csvData), rec.Body.String())
	}
}

func TestExportGradebook_InvalidFormat(t *testing.T) {
	svc := &mockAnalyticsService{}
	mw := &mockAuthMiddleware{userID: 1, userRole: "teacher"}
	h := analyticsHandler.NewAnalyticsHandler(svc, mw)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/courses/13/analytics/export?format=xml", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("courseid", "13")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rec := httptest.NewRecorder()
	h.ExportGradebook(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d", rec.Code)
	}
}

func TestExportGradebook_Unauthorized(t *testing.T) {
	svc := &mockAnalyticsService{}
	mw := &mockAuthMiddleware{userID: 0}
	h := analyticsHandler.NewAnalyticsHandler(svc, mw)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/courses/13/analytics/export", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("courseid", "13")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rec := httptest.NewRecorder()
	h.ExportGradebook(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized, got %d", rec.Code)
	}
}

func TestExportGradebook_Forbidden(t *testing.T) {
	svc := &mockAnalyticsService{
		err: errorsAPP.ErrForbidden,
	}
	mw := &mockAuthMiddleware{userID: 99, userRole: "teacher"}
	h := analyticsHandler.NewAnalyticsHandler(svc, mw)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/courses/13/analytics/export", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("courseid", "13")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rec := httptest.NewRecorder()
	h.ExportGradebook(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden, got %d", rec.Code)
	}
}
