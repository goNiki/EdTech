package auth_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"edtech/internal/domain"
	"edtech/internal/dto"
	authHandler "edtech/internal/interfaces/handlers/auth"
	"edtech/internal/service"
	errorsAPP "edtech/pkg/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockAuthServiceForPreferences struct {
	service.AuthService
	updatePreferencesFunc func(ctx context.Context, userID int64, input domain.UpdatePreferencesInput) (domain.UserPreferences, error)
}

func (m *mockAuthServiceForPreferences) UpdatePreferences(ctx context.Context, userID int64, input domain.UpdatePreferencesInput) (domain.UserPreferences, error) {
	if m.updatePreferencesFunc != nil {
		return m.updatePreferencesFunc(ctx, userID, input)
	}
	return domain.UserPreferences{}, nil
}

type mockAuthMiddlewareForPreferences struct {
	userID int64
}

func (m *mockAuthMiddlewareForPreferences) JWTMiddleware(next http.Handler) http.Handler {
	return next
}

func (m *mockAuthMiddlewareForPreferences) OptionalJWTMiddleware(next http.Handler) http.Handler {
	return next
}

func (m *mockAuthMiddlewareForPreferences) RequireRole(roles []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return next
	}
}

func (m *mockAuthMiddlewareForPreferences) GetUserID(ctx context.Context) int64 {
	return m.userID
}

func (m *mockAuthMiddlewareForPreferences) GetUserRole(ctx context.Context) string {
	return "student"
}

func TestUpdatePreferencesHandler_Success(t *testing.T) {
	mockSvc := &mockAuthServiceForPreferences{
		updatePreferencesFunc: func(ctx context.Context, userID int64, input domain.UpdatePreferencesInput) (domain.UserPreferences, error) {
			assert.Equal(t, int64(10), userID)
			assert.Equal(t, "large", *input.FontScale)
			assert.Equal(t, "wide", *input.ContentWidth)
			assert.Equal(t, "relaxed", *input.LineHeight)
			assert.Equal(t, "sepia", *input.ReadingTheme)
			return domain.UserPreferences{
				FontScale:    "large",
				ContentWidth: "wide",
				LineHeight:   "relaxed",
				ReadingTheme: "sepia",
			}, nil
		},
	}
	mockMW := &mockAuthMiddlewareForPreferences{userID: 10}
	h := authHandler.NewAuthHandler(mockSvc, mockMW)

	body := []byte(`{
		"font_scale": "large",
		"content_width": "wide",
		"line_height": "relaxed",
		"reading_theme": "sepia"
	}`)
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/auth/preferences", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.UpdatePreferences(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp dto.UpdatePreferencesResponse
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, "preferences updated successfully", resp.Message)
	assert.Equal(t, "large", resp.Data.Preferences.FontScale)
	assert.Equal(t, "wide", resp.Data.Preferences.ContentWidth)
	assert.Equal(t, "relaxed", resp.Data.Preferences.LineHeight)
	assert.Equal(t, "sepia", resp.Data.Preferences.ReadingTheme)
}

func TestUpdatePreferencesHandler_InvalidFontScale_Returns400(t *testing.T) {
	mockSvc := &mockAuthServiceForPreferences{}
	mockMW := &mockAuthMiddlewareForPreferences{userID: 10}
	h := authHandler.NewAuthHandler(mockSvc, mockMW)

	body := []byte(`{"font_scale": "giant"}`)
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/auth/preferences", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.UpdatePreferences(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdatePreferencesHandler_Unauthorized_Returns401(t *testing.T) {
	mockSvc := &mockAuthServiceForPreferences{}
	mockMW := &mockAuthMiddlewareForPreferences{userID: 0}
	h := authHandler.NewAuthHandler(mockSvc, mockMW)

	body := []byte(`{"font_scale": "large"}`)
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/auth/preferences", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.UpdatePreferences(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestUpdatePreferencesHandler_ServiceInvalidPreferences_Returns400(t *testing.T) {
	mockSvc := &mockAuthServiceForPreferences{
		updatePreferencesFunc: func(ctx context.Context, userID int64, input domain.UpdatePreferencesInput) (domain.UserPreferences, error) {
			return domain.UserPreferences{}, errorsAPP.ErrInvalidPreferences
		},
	}
	mockMW := &mockAuthMiddlewareForPreferences{userID: 10}
	h := authHandler.NewAuthHandler(mockSvc, mockMW)

	body := []byte(`{"font_scale": "large"}`)
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/auth/preferences", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.UpdatePreferences(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
