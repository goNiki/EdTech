package lesson_test

import (
	"context"
	"testing"

	"edtech/internal/domain"
	"edtech/internal/repository"
	"edtech/internal/service"
	lessonService "edtech/internal/service/lesson"
	errorsAPP "edtech/pkg/errors"
)

type mockCourseRepo struct {
	repository.CourseRepository
	course *domain.Course
	err    error
}

func (m *mockCourseRepo) GetCourseByID(ctx context.Context, id int64) (*domain.Course, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.course, nil
}

type mockLessonRepo struct {
	repository.LessonRepository
	lessonsByID map[int64]*domain.Lesson
	allLessons  []domain.Lesson
	getErr      error
}

func (m *mockLessonRepo) GetLessonByID(ctx context.Context, id int64) (*domain.Lesson, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	l, ok := m.lessonsByID[id]
	if !ok {
		return nil, errorsAPP.ErrLessonNotFound
	}
	return l, nil
}

func (m *mockLessonRepo) GetLessonsByCourseID(ctx context.Context, courseID int64) ([]domain.Lesson, error) {
	return m.allLessons, nil
}

type mockSectionRepo struct {
	repository.SectionRepository
	sections []domain.Section
}

func (m *mockSectionRepo) ListSectionsByCourseID(ctx context.Context, courseID int64) ([]domain.Section, error) {
	return m.sections, nil
}

type mockAccessService struct {
	service.AccessService
	canView bool
	err     error
}

func (m *mockAccessService) CanViewCourse(ctx context.Context, course *domain.Course, userID int64) (bool, error) {
	if m.err != nil {
		return false, m.err
	}
	return m.canView, nil
}

type mockProgressRepo struct {
	repository.ProgressRepository
	progresses []domain.LessonProgress
}

func (m *mockProgressRepo) GetAllLessonProgressByCourse(ctx context.Context, userID, courseID int64) ([]domain.LessonProgress, error) {
	return m.progresses, nil
}

func setupFixture() (*mockCourseRepo, *mockLessonRepo, *mockSectionRepo, *mockAccessService, *mockProgressRepo) {
	courseRepo := &mockCourseRepo{
		course: &domain.Course{
			Id:    10,
			Title: "Тестовый курс Go",
			Slug:  "test-go-course",
		},
	}

	sec1ID := int64(101)
	sec2ID := int64(102)

	sections := []domain.Section{
		{ID: sec1ID, CourseID: 10, Title: "Модуль 1", Position: 1},
		{ID: sec2ID, CourseID: 10, Title: "Модуль 2", Position: 2},
	}
	sectionRepo := &mockSectionRepo{sections: sections}

	l1 := domain.Lesson{ID: 1, CourseID: 10, SectionID: &sec1ID, Title: "Урок 1.1", Position: 1}
	l2 := domain.Lesson{ID: 2, CourseID: 10, SectionID: &sec1ID, Title: "Урок 1.2", Position: 2}
	l3 := domain.Lesson{ID: 3, CourseID: 10, SectionID: &sec2ID, Title: "Урок 2.1", Position: 1}
	l4 := domain.Lesson{ID: 4, CourseID: 10, SectionID: &sec2ID, Title: "Урок 2.2", Position: 2}

	lessonsByID := map[int64]*domain.Lesson{
		1: &l1,
		2: &l2,
		3: &l3,
		4: &l4,
	}
	allLessons := []domain.Lesson{l1, l2, l3, l4}

	lessonRepo := &mockLessonRepo{
		lessonsByID: lessonsByID,
		allLessons:  allLessons,
	}

	accessSvc := &mockAccessService{canView: true}

	score100 := 100
	progresses := []domain.LessonProgress{
		{LessonID: 1, UserID: 1, Status: domain.ProgressStatusCompleted, Score: &score100},
		{LessonID: 2, UserID: 1, Status: domain.ProgressStatusInProgress, Score: nil},
	}
	progressRepo := &mockProgressRepo{progresses: progresses}

	return courseRepo, lessonRepo, sectionRepo, accessSvc, progressRepo
}

func TestGetLessonNavigationContext_MiddleLesson_CrossSectionNavigation(t *testing.T) {
	cRepo, lRepo, sRepo, aSvc, pRepo := setupFixture()
	svc := lessonService.NewLessonService(cRepo, lRepo, sRepo, aSvc, pRepo)

	ctx := context.Background()
	// Current lesson is 2 (last lesson of Section 1)
	nav, err := svc.GetLessonNavigationContext(ctx, 1, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if nav.CurrentLesson.ID != 2 {
		t.Errorf("expected CurrentLesson.ID 2, got %d", nav.CurrentLesson.ID)
	}
	if nav.Course.ID != 10 {
		t.Errorf("expected Course.ID 10, got %d", nav.Course.ID)
	}

	// Prev lesson should be 1
	if nav.PrevLesson == nil || nav.PrevLesson.ID != 1 {
		t.Fatalf("expected PrevLesson ID 1, got %+v", nav.PrevLesson)
	}
	if nav.PrevLesson.Title != "Урок 1.1" {
		t.Errorf("expected PrevLesson.Title 'Урок 1.1', got %s", nav.PrevLesson.Title)
	}

	// Next lesson should cross into Section 2 -> Lesson 3!
	if nav.NextLesson == nil || nav.NextLesson.ID != 3 {
		t.Fatalf("expected NextLesson ID 3, got %+v", nav.NextLesson)
	}
	if nav.NextLesson.Title != "Урок 2.1" {
		t.Errorf("expected NextLesson.Title 'Урок 2.1', got %s", nav.NextLesson.Title)
	}

	// Syllabus checks
	if len(nav.Syllabus) != 2 {
		t.Fatalf("expected 2 syllabus sections, got %d", len(nav.Syllabus))
	}
	sec1 := nav.Syllabus[0]
	if sec1.SectionID != 101 || len(sec1.Lessons) != 2 {
		t.Fatalf("expected section 101 with 2 lessons, got %+v", sec1)
	}
	if !sec1.Lessons[0].IsCompleted || sec1.Lessons[0].Score != 100 {
		t.Errorf("expected lesson 1 to be completed with score 100, got completed=%v score=%d", sec1.Lessons[0].IsCompleted, sec1.Lessons[0].Score)
	}
	if sec1.Lessons[1].IsCompleted {
		t.Errorf("expected lesson 2 to not be completed")
	}
}

func TestGetLessonNavigationContext_FirstLesson_PrevNull(t *testing.T) {
	cRepo, lRepo, sRepo, aSvc, pRepo := setupFixture()
	svc := lessonService.NewLessonService(cRepo, lRepo, sRepo, aSvc, pRepo)

	ctx := context.Background()
	nav, err := svc.GetLessonNavigationContext(ctx, 1, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if nav.PrevLesson != nil {
		t.Errorf("expected PrevLesson to be nil for first lesson, got %+v", nav.PrevLesson)
	}
	if nav.NextLesson == nil || nav.NextLesson.ID != 2 {
		t.Errorf("expected NextLesson to be 2, got %+v", nav.NextLesson)
	}
}

func TestGetLessonNavigationContext_LastLesson_NextNull(t *testing.T) {
	cRepo, lRepo, sRepo, aSvc, pRepo := setupFixture()
	svc := lessonService.NewLessonService(cRepo, lRepo, sRepo, aSvc, pRepo)

	ctx := context.Background()
	nav, err := svc.GetLessonNavigationContext(ctx, 1, 4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if nav.PrevLesson == nil || nav.PrevLesson.ID != 3 {
		t.Errorf("expected PrevLesson to be 3, got %+v", nav.PrevLesson)
	}
	if nav.NextLesson != nil {
		t.Errorf("expected NextLesson to be nil for last lesson, got %+v", nav.NextLesson)
	}
}

func TestGetLessonNavigationContext_SingleLessonCourse(t *testing.T) {
	cRepo, _, sRepo, aSvc, pRepo := setupFixture()
	secID := int64(101)
	single := domain.Lesson{ID: 55, CourseID: 10, SectionID: &secID, Title: "Одиночный урок", Position: 1}
	lRepo := &mockLessonRepo{
		lessonsByID: map[int64]*domain.Lesson{55: &single},
		allLessons:  []domain.Lesson{single},
	}
	svc := lessonService.NewLessonService(cRepo, lRepo, sRepo, aSvc, pRepo)

	ctx := context.Background()
	nav, err := svc.GetLessonNavigationContext(ctx, 1, 55)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if nav.PrevLesson != nil {
		t.Errorf("expected PrevLesson nil, got %+v", nav.PrevLesson)
	}
	if nav.NextLesson != nil {
		t.Errorf("expected NextLesson nil, got %+v", nav.NextLesson)
	}
}

func TestGetLessonNavigationContext_Forbidden(t *testing.T) {
	cRepo, lRepo, sRepo, aSvc, pRepo := setupFixture()
	aSvc.canView = false // Access denied!

	svc := lessonService.NewLessonService(cRepo, lRepo, sRepo, aSvc, pRepo)

	ctx := context.Background()
	_, err := svc.GetLessonNavigationContext(ctx, 999, 1)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetLessonNavigationContext_LessonNotFound(t *testing.T) {
	cRepo, lRepo, sRepo, aSvc, pRepo := setupFixture()
	svc := lessonService.NewLessonService(cRepo, lRepo, sRepo, aSvc, pRepo)

	ctx := context.Background()
	_, err := svc.GetLessonNavigationContext(ctx, 1, 99999)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetLessonNavigationContext_UnsectionedLessons(t *testing.T) {
	cRepo, _, _, aSvc, pRepo := setupFixture()
	secID := int64(101)
	sRepo := &mockSectionRepo{
		sections: []domain.Section{
			{ID: secID, CourseID: 10, Title: "Секция 1", Position: 1},
		},
	}

	l1 := domain.Lesson{ID: 1, CourseID: 10, SectionID: &secID, Title: "Урок в секции", Position: 1}
	l2 := domain.Lesson{ID: 2, CourseID: 10, SectionID: nil, Title: "Общий урок", Position: 1}

	lRepo := &mockLessonRepo{
		lessonsByID: map[int64]*domain.Lesson{1: &l1, 2: &l2},
		allLessons:  []domain.Lesson{l1, l2},
	}

	svc := lessonService.NewLessonService(cRepo, lRepo, sRepo, aSvc, pRepo)

	ctx := context.Background()
	nav, err := svc.GetLessonNavigationContext(ctx, 1, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if nav.NextLesson == nil || nav.NextLesson.ID != 2 {
		t.Fatalf("expected NextLesson ID 2, got %+v", nav.NextLesson)
	}
	if len(nav.Syllabus) != 2 {
		t.Fatalf("expected 2 syllabus sections (1 regular + 1 unsectioned), got %d", len(nav.Syllabus))
	}
	if nav.Syllabus[1].SectionTitle != "Общие уроки" {
		t.Errorf("expected 'Общие уроки', got %s", nav.Syllabus[1].SectionTitle)
	}
}
