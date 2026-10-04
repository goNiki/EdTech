package quiz_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"edtech/internal/domain"
	"edtech/internal/dto"
	quizHandler "edtech/internal/interfaces/handlers/quiz"
	"edtech/internal/interfaces/middleware/auth"
	"edtech/internal/service"
	errorsAPP "edtech/pkg/errors"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockQuizServiceForFeedbackHandler struct {
	service.QuizServices
	feedback *domain.StudentHomeworkFeedback
	err      error
}

func (m *mockQuizServiceForFeedbackHandler) GetStudentHomeworkFeedback(ctx context.Context, userID, lessonID int64) (*domain.StudentHomeworkFeedback, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.feedback, nil
}

type mockAuthMiddleware struct {
	auth.AuthMiddleware
	userID int64
}

func (m *mockAuthMiddleware) GetUserID(ctx context.Context) int64 {
	return m.userID
}

func TestGetStudentHomeworkFeedback_Success(t *testing.T) {
	now := time.Now()
	attemptID := int64(3)
	avatar := "/static/uploads/avatars/teacher.png"
	feedbackText := "Отличная аргументация тезиса."
	isCorrect := true

	mockSvc := &mockQuizServiceForFeedbackHandler{
		feedback: &domain.StudentHomeworkFeedback{
			HasSubmission: true,
			Status:        "graded",
			AttemptID:     &attemptID,
			SubmittedAt:   &now,
			GradedAt:      &now,
			Teacher: &domain.HomeworkTeacherInfo{
				ID:        2,
				Name:      "Никита Преподаватель",
				AvatarURL: &avatar,
			},
			Answers: []domain.HomeworkAnswerDetail{
				{
					AnswerID:      3,
					QuestionText:  "Напишите развернутое сочинение-рассуждение...",
					StudentAnswer: "В данном тексте автор поднимает важную проблему...",
					Points:        23,
					MaxPoints:     25,
					IsCorrect:     &isCorrect,
					Feedback:      &feedbackText,
				},
			},
		},
	}

	authMw := &mockAuthMiddleware{userID: 5}
	h := quizHandler.NewQuizHandler(mockSvc, authMw)

	r := chi.NewRouter()
	r.Get("/api/v1/lessons/{id}/homework-feedback", h.GetStudentHomeworkFeedback)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lessons/10/homework-feedback", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var resp dto.HomeworkFeedbackResponse
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, "success", resp.Message)
	assert.True(t, resp.Data.HasSubmission)
	assert.Equal(t, "graded", resp.Data.Status)
	assert.Equal(t, attemptID, *resp.Data.AttemptID)
	assert.NotNil(t, resp.Data.Teacher)
	assert.Equal(t, int64(2), resp.Data.Teacher.ID)
	assert.Equal(t, "Никита Преподаватель", resp.Data.Teacher.Name)
	assert.Len(t, resp.Data.Answers, 1)
	assert.Equal(t, 23, resp.Data.Answers[0].Points)
	assert.Equal(t, 25, resp.Data.Answers[0].MaxPoints)
	assert.NotNil(t, resp.Data.Answers[0].IsCorrect)
	assert.True(t, *resp.Data.Answers[0].IsCorrect)
	assert.Equal(t, feedbackText, *resp.Data.Answers[0].Feedback)
}

func TestGetStudentHomeworkFeedback_Unauthorized(t *testing.T) {
	mockSvc := &mockQuizServiceForFeedbackHandler{}
	authMw := &mockAuthMiddleware{userID: 0}
	h := quizHandler.NewQuizHandler(mockSvc, authMw)

	r := chi.NewRouter()
	r.Get("/api/v1/lessons/{id}/homework-feedback", h.GetStudentHomeworkFeedback)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lessons/10/homework-feedback", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestGetStudentHomeworkFeedback_InvalidID(t *testing.T) {
	mockSvc := &mockQuizServiceForFeedbackHandler{}
	authMw := &mockAuthMiddleware{userID: 5}
	h := quizHandler.NewQuizHandler(mockSvc, authMw)

	r := chi.NewRouter()
	r.Get("/api/v1/lessons/{id}/homework-feedback", h.GetStudentHomeworkFeedback)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lessons/invalid/homework-feedback", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetStudentHomeworkFeedback_Forbidden(t *testing.T) {
	mockSvc := &mockQuizServiceForFeedbackHandler{
		err: errorsAPP.ErrForbidden,
	}
	authMw := &mockAuthMiddleware{userID: 5}
	h := quizHandler.NewQuizHandler(mockSvc, authMw)

	r := chi.NewRouter()
	r.Get("/api/v1/lessons/{id}/homework-feedback", h.GetStudentHomeworkFeedback)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lessons/10/homework-feedback", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestGetStudentHomeworkFeedback_NotFound(t *testing.T) {
	mockSvc := &mockQuizServiceForFeedbackHandler{
		err: errorsAPP.ErrLessonNotFound,
	}
	authMw := &mockAuthMiddleware{userID: 5}
	h := quizHandler.NewQuizHandler(mockSvc, authMw)

	r := chi.NewRouter()
	r.Get("/api/v1/lessons/{id}/homework-feedback", h.GetStudentHomeworkFeedback)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lessons/999/homework-feedback", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}
