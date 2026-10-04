package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"edtech/internal/domain"
	"edtech/internal/repository"
	authService "edtech/internal/service/auth"
	errorsAPP "edtech/pkg/errors"
)

type mockAdminUserRepo struct {
	repository.UserRepository
	users        map[int64]*domain.User
	getUserErr   error
	listUsersErr error
}

func newMockAdminUserRepo() *mockAdminUserRepo {
	return &mockAdminUserRepo{
		users: make(map[int64]*domain.User),
	}
}

func (m *mockAdminUserRepo) GetUserByID(ctx context.Context, userID int64) (*domain.User, error) {
	if m.getUserErr != nil {
		return nil, m.getUserErr
	}
	u, ok := m.users[userID]
	if !ok {
		return nil, errorsAPP.ErrUserNotFound
	}
	return u, nil
}

func (m *mockAdminUserRepo) ListUsers(ctx context.Context, filter domain.UserFilter, pagination domain.Pagination) ([]domain.User, int64, error) {
	if m.listUsersErr != nil {
		return nil, 0, m.listUsersErr
	}

	result := make([]domain.User, 0)
	for _, u := range m.users {
		if filter.Role != nil && u.Role != *filter.Role {
			continue
		}
		if filter.IsBanned != nil && u.IsBanned != *filter.IsBanned {
			continue
		}
		result = append(result, *u)
	}

	return result, int64(len(result)), nil
}

func TestListUsersForAdmin_Success(t *testing.T) {
	repo := newMockAdminUserRepo()
	repo.users[1] = &domain.User{ID: 1, Role: domain.RoleAdmin}
	repo.users[2] = &domain.User{ID: 2, Role: domain.RoleStudent, CreatedAt: time.Now()}
	repo.users[3] = &domain.User{ID: 3, Role: domain.RoleTeacher, CreatedAt: time.Now()}

	svc := authService.NewAuthService(repo, nil, nil, nil, nil)

	users, total, err := svc.ListUsersForAdmin(context.Background(), 1, domain.UserFilter{}, domain.Pagination{Page: 1, PageSize: 20})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if total != 3 {
		t.Errorf("expected 3 users total, got %d", total)
	}
	if len(users) != 3 {
		t.Errorf("expected 3 users returned, got %d", len(users))
	}
}

func TestListUsersForAdmin_FilterByRole(t *testing.T) {
	repo := newMockAdminUserRepo()
	repo.users[1] = &domain.User{ID: 1, Role: domain.RoleAdmin}
	repo.users[2] = &domain.User{ID: 2, Role: domain.RoleStudent}
	repo.users[3] = &domain.User{ID: 3, Role: domain.RoleTeacher}

	svc := authService.NewAuthService(repo, nil, nil, nil, nil)

	teacherRole := domain.RoleTeacher
	users, total, err := svc.ListUsersForAdmin(context.Background(), 1, domain.UserFilter{Role: &teacherRole}, domain.Pagination{Page: 1, PageSize: 20})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if total != 1 {
		t.Errorf("expected 1 user, got %d", total)
	}
	if len(users) != 1 || users[0].ID != 3 {
		t.Errorf("expected teacher user with ID 3, got %+v", users)
	}
}

func TestListUsersForAdmin_ForbiddenForNonAdmin(t *testing.T) {
	repo := newMockAdminUserRepo()
	repo.users[2] = &domain.User{ID: 2, Role: domain.RoleStudent}

	svc := authService.NewAuthService(repo, nil, nil, nil, nil)

	_, _, err := svc.ListUsersForAdmin(context.Background(), 2, domain.UserFilter{}, domain.Pagination{Page: 1, PageSize: 20})
	if !errors.Is(err, errorsAPP.ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got: %v", err)
	}
}

func TestListUsersForAdmin_AdminNotFound(t *testing.T) {
	repo := newMockAdminUserRepo()
	svc := authService.NewAuthService(repo, nil, nil, nil, nil)

	_, _, err := svc.ListUsersForAdmin(context.Background(), 999, domain.UserFilter{}, domain.Pagination{Page: 1, PageSize: 20})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
