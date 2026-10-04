package quiz_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"edtech/internal/domain"
	quizHandler "edtech/internal/interfaces/handlers/quiz"
	"edtech/internal/service"
	errorsAPP "edtech/pkg/errors"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockQuizServiceForGradingHandler struct {
	service.QuizServices
	attempt *domain.QuizAttempt
	err     error
}

func (m *mockQuizServiceForGradingHandler) GradeAttemptAnswer(ctx context.Context, teacherID, attemptID, answerID int64, points int, feedback *string) (*domain.QuizAttempt, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.attempt, nil
}

func TestGradeAttemptAnswerHandler_NegativePoints_ValidationFailed(t *testing.T) {
	mockSvc := &mockQuizServiceForGradingHandler{}
	mockMW := &mockAuthMiddleware{userID: 1}

	h := quizHandler.NewQuizHandler(mockSvc, mockMW)

	body := []byte(`{"points": -1}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/quizzes/attempts/18/answers/16/grade", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("attempt_id", "18")
	rctx.URLParams.Add("answer_id", "16")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rec := httptest.NewRecorder()
	h.GradeAttemptAnswer(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGradeAttemptAnswerHandler_PointsExceedMax_Returns400(t *testing.T) {
	mockSvc := &mockQuizServiceForGradingHandler{
		err: fmt.Errorf("%w: балл не может быть меньше 0 или превышать максимальный балл задания (25)", errorsAPP.ErrInvalidGradePoints),
	}
	mockMW := &mockAuthMiddleware{userID: 1}

	h := quizHandler.NewQuizHandler(mockSvc, mockMW)

	body := []byte(`{"points": 999, "feedback": "Over max"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/quizzes/attempts/18/answers/16/grade", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("attempt_id", "18")
	rctx.URLParams.Add("answer_id", "16")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rec := httptest.NewRecorder()
	h.GradeAttemptAnswer(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var errResp map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &errResp)
	require.NoError(t, err)
	assert.Equal(t, "INVALID_GRADE_POINTS", errResp["code"])
	assert.Contains(t, errResp["error"], "максимальный балл задания (25)")
}

func TestGradeAttemptAnswerHandler_ValidPoints_Success(t *testing.T) {
	mockSvc := &mockQuizServiceForGradingHandler{
		attempt: &domain.QuizAttempt{
			ID:     18,
			QuizID: 10,
			UserID: 5,
			Score:  80,
			Passed: true,
		},
	}
	mockMW := &mockAuthMiddleware{userID: 1}

	h := quizHandler.NewQuizHandler(mockSvc, mockMW)

	body := []byte(`{"points": 20, "feedback": "Отлично"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/quizzes/attempts/18/answers/16/grade", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("attempt_id", "18")
	rctx.URLParams.Add("answer_id", "16")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rec := httptest.NewRecorder()
	h.GradeAttemptAnswer(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}
