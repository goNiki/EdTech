package auth_test

import (
	"context"
	"testing"

	"edtech/internal/domain"
	"edtech/internal/repository"
	authService "edtech/internal/service/auth"
)

type mockUpdateProfileUserRepo struct {
	repository.UserRepository
	user           *domain.User
	getUserErr     error
	updatedUser    *domain.User
	updateProfileErr error
}

func (m *mockUpdateProfileUserRepo) GetUserByID(ctx context.Context, userID int64) (*domain.User, error) {
	if m.getUserErr != nil {
		return nil, m.getUserErr
	}
	return m.user, nil
}

func (m *mockUpdateProfileUserRepo) UpdateProfile(ctx context.Context, user *domain.User) error {
	if m.updateProfileErr != nil {
		return m.updateProfileErr
	}
	m.updatedUser = user
	return nil
}

func TestUpdateProfile_Headline(t *testing.T) {
	headlineOld := "Старый заголовок"
	headlineNew := "Ведущий преподаватель математики"
	bioNew := "Новая биография"

	user := &domain.User{
		ID:        42,
		Email:     "teacher@example.com",
		Username:  "teacher_ivan",
		Headline:  &headlineOld,
		Role:      domain.RoleTeacher,
		IsActive:  true,
	}

	mockRepo := &mockUpdateProfileUserRepo{user: user}
	svc := authService.NewAuthService(mockRepo, nil, nil, nil, nil)

	input := domain.UpdateProfileInput{
		Headline: &headlineNew,
		Bio:      &bioNew,
	}

	updated, err := svc.UpdateProfile(context.Background(), 42, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if updated == nil {
		t.Fatal("expected updated user, got nil")
	}

	if updated.Headline == nil || *updated.Headline != headlineNew {
		t.Errorf("expected headline %q, got %v", headlineNew, updated.Headline)
	}

	if updated.Bio == nil || *updated.Bio != bioNew {
		t.Errorf("expected bio %q, got %v", bioNew, updated.Bio)
	}

	if mockRepo.updatedUser == nil {
		t.Fatal("expected repo.UpdateProfile to be called")
	}

	if mockRepo.updatedUser.Headline == nil || *mockRepo.updatedUser.Headline != headlineNew {
		t.Errorf("expected repo headline %q, got %v", headlineNew, mockRepo.updatedUser.Headline)
	}
}
