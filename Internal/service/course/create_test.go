package course_test

import (
	"context"
	"fmt"
	"testing"

	"edtech/internal/domain"
	"edtech/internal/repository"
	"edtech/internal/service"
	courseService "edtech/internal/service/course"

	"github.com/jackc/pgx/v5"
)

type mockTxManager struct{}

func (m *mockTxManager) WithTX(ctx context.Context, opts pgx.TxOptions, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

type mockCourseRepo struct {
	repository.CourseRepository
	coursesBySlug map[string]*domain.Course
	nextID        int64
}

func newMockCourseRepo() *mockCourseRepo {
	return &mockCourseRepo{
		coursesBySlug: make(map[string]*domain.Course),
		nextID:        1,
	}
}

func (m *mockCourseRepo) ExistsBySlug(ctx context.Context, slug string) (bool, error) {
	_, ok := m.coursesBySlug[slug]
	return ok, nil
}

func (m *mockCourseRepo) CreateCourse(ctx context.Context, c *domain.Course) (*domain.Course, error) {
	if _, exists := m.coursesBySlug[c.Slug]; exists {
		return nil, fmt.Errorf("duplicate key value violates unique constraint slug %s", c.Slug)
	}
	c.Id = m.nextID
	m.nextID++
	m.coursesBySlug[c.Slug] = c
	return c, nil
}

type mockEnrolledRepo struct {
	repository.EnrolledRepository
}

func (m *mockEnrolledRepo) EnrollUserToCourse(ctx context.Context, e domain.EnrolledInCourse) error {
	return nil
}

func TestCreateCourse_SlugCollisionResolution(t *testing.T) {
	courseRepo := newMockCourseRepo()
	enrolledRepo := &mockEnrolledRepo{}
	txMgr := &mockTxManager{}

	svc := courseService.NewCourseService(courseRepo, nil, nil, nil, enrolledRepo, txMgr)

	// Course 1
	c1 := &domain.Course{
		Title:      "Golang Basics",
		Slug:       "golang-basics",
		CreatedBy:  1,
		Visibility: domain.VisibilityPublic,
	}
	created1, err := svc.CreateCourse(context.Background(), c1)
	if err != nil {
		t.Fatalf("unexpected error creating course 1: %v", err)
	}
	if created1.Slug != "golang-basics" {
		t.Errorf("expected slug 'golang-basics', got %q", created1.Slug)
	}

	// Course 2 with identical slug
	c2 := &domain.Course{
		Title:      "Golang Basics Duplicate",
		Slug:       "golang-basics",
		CreatedBy:  1,
		Visibility: domain.VisibilityPublic,
	}
	created2, err := svc.CreateCourse(context.Background(), c2)
	if err != nil {
		t.Fatalf("unexpected error creating course 2: %v", err)
	}
	if created2.Slug != "golang-basics-1" {
		t.Errorf("expected auto-resolved slug 'golang-basics-1', got %q", created2.Slug)
	}

	// Course 3 with identical slug
	c3 := &domain.Course{
		Title:      "Golang Basics Third",
		Slug:       "golang-basics",
		CreatedBy:  1,
		Visibility: domain.VisibilityPublic,
	}
	created3, err := svc.CreateCourse(context.Background(), c3)
	if err != nil {
		t.Fatalf("unexpected error creating course 3: %v", err)
	}
	if created3.Slug != "golang-basics-2" {
		t.Errorf("expected auto-resolved slug 'golang-basics-2', got %q", created3.Slug)
	}

	// Course 4 with unnormalized dirty slug
	c4 := &domain.Course{
		Title:      "Advanced Go",
		Slug:       "  Advanced_Go & Concurrency!!  ",
		CreatedBy:  1,
		Visibility: domain.VisibilityPublic,
	}
	created4, err := svc.CreateCourse(context.Background(), c4)
	if err != nil {
		t.Fatalf("unexpected error creating course 4: %v", err)
	}
	if created4.Slug != "advanced-go-concurrency" {
		t.Errorf("expected normalized slug 'advanced-go-concurrency', got %q", created4.Slug)
	}
}

var _ service.CourseServices = (service.CourseServices)(nil)
