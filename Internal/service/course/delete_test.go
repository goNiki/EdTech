package course_test

import (
	"context"
	"errors"
	"testing"

	"edtech/internal/domain"
	"edtech/internal/repository"
	"edtech/internal/service"
	courseService "edtech/internal/service/course"
	errorsAPP "edtech/pkg/errors"
)

type mockAccessService struct {
	service.AccessService
	canDelete bool
	err       error
}

func (m *mockAccessService) CanDeleteCourse(ctx context.Context, course *domain.Course, userID int64) (bool, error) {
	if m.err != nil {
		return false, m.err
	}
	return m.canDelete, nil
}

type mockDeleteCourseRepo struct {
	repository.CourseRepository
	course    *domain.Course
	getErr    error
	deletedID int64
	delErr    error
}

func (m *mockDeleteCourseRepo) GetCourseByID(ctx context.Context, id int64) (*domain.Course, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	return m.course, nil
}

func (m *mockDeleteCourseRepo) DeleteCourse(ctx context.Context, courseID int64) error {
	if m.delErr != nil {
		return m.delErr
	}
	m.deletedID = courseID
	return nil
}

func TestDeleteCourse_ForbiddenForNonCreator(t *testing.T) {
	ctx := context.Background()
	course := &domain.Course{Id: 10, CreatedBy: 1}
	repo := &mockDeleteCourseRepo{course: course}
	access := &mockAccessService{canDelete: false}

	svc := courseService.NewCourseService(repo, nil, nil, access, nil, nil)

	err := svc.DeleteCourse(ctx, 10, 2)
	if !errors.Is(err, errorsAPP.ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestDeleteCourse_SuccessForCreator(t *testing.T) {
	ctx := context.Background()
	course := &domain.Course{Id: 10, CreatedBy: 1}
	repo := &mockDeleteCourseRepo{course: course}
	access := &mockAccessService{canDelete: true}

	svc := courseService.NewCourseService(repo, nil, nil, access, nil, nil)

	err := svc.DeleteCourse(ctx, 10, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repo.deletedID != 10 {
		t.Errorf("expected deletedID 10, got %d", repo.deletedID)
	}
}

func TestDeleteCourse_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := &mockDeleteCourseRepo{getErr: errorsAPP.ErrCourseNotFound}

	svc := courseService.NewCourseService(repo, nil, nil, nil, nil, nil)

	err := svc.DeleteCourse(ctx, 999, 1)
	if !errors.Is(err, errorsAPP.ErrCourseNotFound) {
		t.Fatalf("expected ErrCourseNotFound, got %v", err)
	}
}
