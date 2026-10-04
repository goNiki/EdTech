package auth_test

import (
	"context"
	"errors"
	"testing"

	"edtech/internal/domain"
	"edtech/internal/repository"
	authService "edtech/internal/service/auth"
	errorsAPP "edtech/pkg/errors"
)

type mockPreferencesUserRepo struct {
	repository.UserRepository
	updatePreferencesFunc func(ctx context.Context, userID int64, input domain.UpdatePreferencesInput) (domain.UserPreferences, error)
}

func (m *mockPreferencesUserRepo) UpdatePreferences(ctx context.Context, userID int64, input domain.UpdatePreferencesInput) (domain.UserPreferences, error) {
	if m.updatePreferencesFunc != nil {
		return m.updatePreferencesFunc(ctx, userID, input)
	}
	return domain.UserPreferences{}, nil
}

func TestUpdatePreferences_Success(t *testing.T) {
	fontScale := "large"
	theme := "sepia"
	width := "wide"
	lineHeight := "relaxed"

	mockRepo := &mockPreferencesUserRepo{
		updatePreferencesFunc: func(ctx context.Context, userID int64, input domain.UpdatePreferencesInput) (domain.UserPreferences, error) {
			if userID != 10 {
				t.Fatalf("expected userID 10, got %d", userID)
			}
			return domain.UserPreferences{
				FontScale:    *input.FontScale,
				ContentWidth: *input.ContentWidth,
				LineHeight:   *input.LineHeight,
				ReadingTheme: *input.ReadingTheme,
			}, nil
		},
	}

	svc := authService.NewAuthService(mockRepo, nil, nil, nil, nil)
	input := domain.UpdatePreferencesInput{
		FontScale:    &fontScale,
		ContentWidth: &width,
		LineHeight:   &lineHeight,
		ReadingTheme: &theme,
	}

	res, err := svc.UpdatePreferences(context.Background(), 10, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.FontScale != "large" || res.ReadingTheme != "sepia" || res.ContentWidth != "wide" || res.LineHeight != "relaxed" {
		t.Fatalf("unexpected preferences returned: %+v", res)
	}
}

func TestUpdatePreferences_ValidationFailure(t *testing.T) {
	invalidScale := "super-huge"
	svc := authService.NewAuthService(&mockPreferencesUserRepo{}, nil, nil, nil, nil)

	input := domain.UpdatePreferencesInput{
		FontScale: &invalidScale,
	}

	_, err := svc.UpdatePreferences(context.Background(), 1, input)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, errorsAPP.ErrInvalidPreferences) {
		t.Fatalf("expected ErrInvalidPreferences, got %v", err)
	}
}

func TestUpdatePreferences_UserNotFound(t *testing.T) {
	mockRepo := &mockPreferencesUserRepo{
		updatePreferencesFunc: func(ctx context.Context, userID int64, input domain.UpdatePreferencesInput) (domain.UserPreferences, error) {
			return domain.UserPreferences{}, errorsAPP.ErrUserNotFound
		},
	}

	svc := authService.NewAuthService(mockRepo, nil, nil, nil, nil)
	fontScale := "medium"
	_, err := svc.UpdatePreferences(context.Background(), 999, domain.UpdatePreferencesInput{FontScale: &fontScale})
	if !errors.Is(err, errorsAPP.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}
