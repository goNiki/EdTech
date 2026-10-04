package enrollment_test

import (
	"context"
	"errors"
	"testing"

	"edtech/internal/domain"
	"edtech/internal/repository"
	"edtech/internal/service"
	enrollmentService "edtech/internal/service/enrollment"
	errorsAPP "edtech/pkg/errors"

	"github.com/jackc/pgx/v5"
)

type mockTxManager struct{}

func (m *mockTxManager) WithTX(ctx context.Context, opts pgx.TxOptions, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

type mockEnrollmentRepo struct {
	repository.EnrolledRepository
	roles        map[int64]string // userID -> role
	unenrolledID int64
}

func (m *mockEnrollmentRepo) GetRoleUserInCourse(ctx context.Context, userID int64, courseID int64) (string, error) {
	role, ok := m.roles[userID]
	if !ok {
		return "", errorsAPP.ErrNotEnrolled
	}
	return role, nil
}

func (m *mockEnrollmentRepo) UnenrollUser(ctx context.Context, userID int64, courseID int64) error {
	m.unenrolledID = userID
	delete(m.roles, userID)
	return nil
}

type mockCourseRepo struct {
	repository.CourseRepository
	enrolledCount int
	course        *domain.Course
}

func (m *mockCourseRepo) GetCourseByID(ctx context.Context, id int64) (*domain.Course, error) {
	if m.course != nil {
		return m.course, nil
	}
	return &domain.Course{Id: id, CreatedBy: 1}, nil
}

func (m *mockCourseRepo) DecrementEnrolledCount(ctx context.Context, courseID int64) error {
	if m.enrolledCount > 0 {
		m.enrolledCount--
	}
	return nil
}

type mockAccessService struct {
	service.AccessService
	canManage bool
}

func (m *mockAccessService) CanManageCourseUsers(ctx context.Context, course *domain.Course, userID int64) (bool, error) {
	return m.canManage, nil
}

func TestUnenrollUser_CreatorCannotUnenroll(t *testing.T) {
	enrolledRepo := &mockEnrollmentRepo{
		roles: map[int64]string{1: "creator"},
	}
	courseRepo := &mockCourseRepo{enrolledCount: 5}
	txMgr := &mockTxManager{}

	svc := enrollmentService.NewEnrolmentService(courseRepo, enrolledRepo, nil, nil, txMgr)

	err := svc.UnenrollUser(context.Background(), 1, 100)
	if err == nil {
		t.Fatal("expected error when unenrolling creator, got nil")
	}
	if !errors.Is(err, errorsAPP.ErrCreatorCannotUnenroll) {
		t.Fatalf("expected ErrCreatorCannotUnenroll, got %v", err)
	}
}

func TestUnenrollUser_StudentSuccess(t *testing.T) {
	enrolledRepo := &mockEnrollmentRepo{
		roles: map[int64]string{2: "student"},
	}
	courseRepo := &mockCourseRepo{enrolledCount: 5}
	txMgr := &mockTxManager{}

	svc := enrollmentService.NewEnrolmentService(courseRepo, enrolledRepo, nil, nil, txMgr)

	err := svc.UnenrollUser(context.Background(), 2, 100)
	if err != nil {
		t.Fatalf("unexpected error unenrolling student: %v", err)
	}

	if enrolledRepo.unenrolledID != 2 {
		t.Errorf("expected unenrolled user 2, got %d", enrolledRepo.unenrolledID)
	}
	if courseRepo.enrolledCount != 4 {
		t.Errorf("expected enrolled count 4, got %d", courseRepo.enrolledCount)
	}
}

func TestTeacherUnenrollUser_Forbidden(t *testing.T) {
	enrolledRepo := &mockEnrollmentRepo{
		roles: map[int64]string{2: "student"},
	}
	courseRepo := &mockCourseRepo{enrolledCount: 5}
	accessSvc := &mockAccessService{canManage: false}
	txMgr := &mockTxManager{}

	svc := enrollmentService.NewEnrolmentService(courseRepo, enrolledRepo, nil, accessSvc, txMgr)

	err := svc.TeacherUnenrollUser(context.Background(), 999, 2, 100)
	if err == nil {
		t.Fatal("expected forbidden error, got nil")
	}
	if !errors.Is(err, errorsAPP.ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestTeacherUnenrollUser_Success(t *testing.T) {
	enrolledRepo := &mockEnrollmentRepo{
		roles: map[int64]string{2: "student"},
	}
	courseRepo := &mockCourseRepo{enrolledCount: 5}
	accessSvc := &mockAccessService{canManage: true}
	txMgr := &mockTxManager{}

	svc := enrollmentService.NewEnrolmentService(courseRepo, enrolledRepo, nil, accessSvc, txMgr)

	err := svc.TeacherUnenrollUser(context.Background(), 1, 2, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if enrolledRepo.unenrolledID != 2 {
		t.Errorf("expected unenrolled student 2, got %d", enrolledRepo.unenrolledID)
	}
}
