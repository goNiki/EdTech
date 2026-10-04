package auth_test

import (
	"context"
	"errors"
	"testing"

	"edtech/internal/domain"
	"edtech/internal/repository"
	authService "edtech/internal/service/auth"
	errorsAPP "edtech/pkg/errors"

	"github.com/jackc/pgx/v5"
)

type mockTxManager struct{}

func (m *mockTxManager) WithTX(ctx context.Context, opts pgx.TxOptions, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

type mockUserRepo struct {
	repository.UserRepository
	user          *domain.User
	getUserErr    error
	updatedPass   string
	updatedUserID int64
	updatePassErr error
}

func (m *mockUserRepo) GetUserByID(ctx context.Context, userID int64) (*domain.User, error) {
	if m.getUserErr != nil {
		return nil, m.getUserErr
	}
	return m.user, nil
}

func (m *mockUserRepo) UpdatePassword(ctx context.Context, userID int64, passHash string) error {
	if m.updatePassErr != nil {
		return m.updatePassErr
	}
	m.updatedUserID = userID
	m.updatedPass = passHash
	return nil
}

type mockRefreshRepo struct {
	repository.RefreshRepository
	revokedUserID int64
	revokeErr     error
}

func (m *mockRefreshRepo) DeleteAllByUserID(ctx context.Context, userID int64) error {
	if m.revokeErr != nil {
		return m.revokeErr
	}
	m.revokedUserID = userID
	return nil
}

type mockHasher struct {
	checkResult bool
	checkErr    error
	hashResult  string
	hashErr     error
}

func (m *mockHasher) Hash(password string) (string, error) {
	if m.hashErr != nil {
		return "", m.hashErr
	}
	if m.hashResult != "" {
		return m.hashResult, nil
	}
	return "hashed_" + password, nil
}

func (m *mockHasher) CheckPassword(passhash string, password string) (bool, error) {
	if m.checkErr != nil {
		return false, m.checkErr
	}
	return m.checkResult, nil
}

func (m *mockHasher) HashRefreshToken(token string) string {
	return "hashed_rt_" + token
}

func TestChangePassword_SamePassword_ReturnsErrSamePassword(t *testing.T) {
	svc := authService.NewAuthService(nil, nil, nil, nil, nil)
	err := svc.ChangePassword(context.Background(), 1, "Password123!", "Password123!")
	if !errors.Is(err, errorsAPP.ErrSamePassword) {
		t.Fatalf("expected ErrSamePassword, got %v", err)
	}
}

func TestChangePassword_ShortPassword_ReturnsErrPasswordTooShort(t *testing.T) {
	svc := authService.NewAuthService(nil, nil, nil, nil, nil)
	err := svc.ChangePassword(context.Background(), 1, "OldPassword123!", "Short1!")
	if !errors.Is(err, errorsAPP.ErrPasswordTooShort) {
		t.Fatalf("expected ErrPasswordTooShort, got %v", err)
	}
}

func TestChangePassword_UserNotFound_ReturnsError(t *testing.T) {
	userRepo := &mockUserRepo{
		getUserErr: errorsAPP.ErrUserNotFound,
	}
	svc := authService.NewAuthService(userRepo, nil, nil, nil, nil)
	err := svc.ChangePassword(context.Background(), 99, "OldPassword123!", "NewStrongPassword456!")
	if err == nil || !errors.Is(err, errorsAPP.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestChangePassword_WrongOldPassword_ReturnsErrInvalidCredentials(t *testing.T) {
	userRepo := &mockUserRepo{
		user: &domain.User{ID: 1, PasswordHash: "stored_hash"},
	}
	hasher := &mockHasher{
		checkResult: false,
	}
	svc := authService.NewAuthService(userRepo, nil, hasher, nil, nil)
	err := svc.ChangePassword(context.Background(), 1, "WrongPassword123!", "NewStrongPassword456!")
	if !errors.Is(err, errorsAPP.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestChangePassword_Success_UpdatesHashAndRevokesTokens(t *testing.T) {
	ctx := context.Background()
	userRepo := &mockUserRepo{
		user: &domain.User{ID: 42, PasswordHash: "stored_hash"},
	}
	refreshRepo := &mockRefreshRepo{}
	hasher := &mockHasher{
		checkResult: true,
		hashResult:  "new_secure_bcrypt_hash",
	}
	txMgr := &mockTxManager{}

	svc := authService.NewAuthService(userRepo, nil, hasher, refreshRepo, txMgr)

	err := svc.ChangePassword(ctx, 42, "OldCorrectPassword123!", "NewStrongPassword456!")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if userRepo.updatedUserID != 42 {
		t.Errorf("expected updated user ID 42, got %d", userRepo.updatedUserID)
	}
	if userRepo.updatedPass != "new_secure_bcrypt_hash" {
		t.Errorf("expected updated password hash 'new_secure_bcrypt_hash', got '%s'", userRepo.updatedPass)
	}
	if refreshRepo.revokedUserID != 42 {
		t.Errorf("expected revoked refresh tokens for user ID 42, got %d", refreshRepo.revokedUserID)
	}
}
