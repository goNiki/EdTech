package analytics_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"edtech/internal/domain"
	analyticsService "edtech/internal/service/analytics"
	errorsAPP "edtech/pkg/errors"
)

type mockAnalyticsRepo struct {
	items   []domain.PendingHomeworkItem
	total   int64
	summary []domain.CoursePendingSummaryItem
	err     error
}

func (m *mockAnalyticsRepo) GetCourseAnalyticsSummary(ctx context.Context, courseID int64) (domain.CourseAnalyticsSummary, error) {
	return domain.CourseAnalyticsSummary{}, nil
}

func (m *mockAnalyticsRepo) ListPendingHomeworks(ctx context.Context, courseID int64, limit, offset int64) ([]domain.PendingHomeworkItem, int64, error) {
	if m.err != nil {
		return nil, 0, m.err
	}
	return m.items, m.total, nil
}

func (m *mockAnalyticsRepo) ListTeacherPendingHomeworks(ctx context.Context, teacherID int64, courseID int64, limit, offset int64) ([]domain.PendingHomeworkItem, int64, []domain.CoursePendingSummaryItem, error) {
	if m.err != nil {
		return nil, 0, nil, m.err
	}
	return m.items, m.total, m.summary, nil
}

func (m *mockAnalyticsRepo) GetStudentDrilldown(ctx context.Context, userID, courseID int64) (*domain.StudentDrilldownReport, error) {
	return nil, nil
}

type mockCourseRepo struct {
	course *domain.Course
	err    error
}

func (m *mockCourseRepo) GetCourseByID(ctx context.Context, id int64) (*domain.Course, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.course, nil
}
func (m *mockCourseRepo) GetCourseBySlug(ctx context.Context, slug string) (*domain.Course, error) {
	return nil, nil
}
func (m *mockCourseRepo) CreateCourse(ctx context.Context, course *domain.Course) (*domain.Course, error) {
	return nil, nil
}
func (m *mockCourseRepo) ArchiveCourse(ctx context.Context, course *domain.Course) error {
	return nil
}
func (m *mockCourseRepo) PublishCourse(ctx context.Context, course *domain.Course) error {
	return nil
}
func (m *mockCourseRepo) UpdateCourseStatus(ctx context.Context, courseID int64, status string) error {
	return nil
}
func (m *mockCourseRepo) UpdateCourse(ctx context.Context, course *domain.Course) error {
	return nil
}
func (m *mockCourseRepo) ListPublicCourses(ctx context.Context, pagination domain.Pagination, filter domain.CourseFilter) ([]domain.Course, error) {
	return nil, nil
}
func (m *mockCourseRepo) ListEnrolledCourses(ctx context.Context, input *domain.InputListMyCourse) ([]domain.Course, error) {
	return nil, nil
}
func (m *mockCourseRepo) CountEnrolledCourses(ctx context.Context, input *domain.InputListMyCourse) (int64, error) {
	return 0, nil
}
func (m *mockCourseRepo) CountCourses(ctx context.Context, filter domain.CourseFilter) (int64, error) {
	return 0, nil
}
func (m *mockCourseRepo) DeleteCourse(ctx context.Context, courseID int64) error {
	return nil
}
func (m *mockCourseRepo) ExistingBySlug(ctx context.Context, slug string) (bool, error) {
	return false, nil
}
func (m *mockCourseRepo) ExistsBySlug(ctx context.Context, slug string) (bool, error) {
	return false, nil
}
func (m *mockCourseRepo) IncrementEnrolledCount(ctx context.Context, courseID int64) error {
	return nil
}
func (m *mockCourseRepo) DecrementEnrolledCount(ctx context.Context, courseID int64) error {
	return nil
}
func (m *mockCourseRepo) UpdateCourseRatingStats(ctx context.Context, courseID int64, rating float64, reviewsCount int) error {
	return nil
}

type mockAccessService struct {
	canEdit bool
	err     error
}

func (m *mockAccessService) BuildCoursePermissions(ctx context.Context, course *domain.Course, userID int64) (domain.CoursePermissions, error) {
	return domain.CoursePermissions{}, nil
}
func (m *mockAccessService) CanAccessCourseObject(ctx context.Context, course *domain.Course, userID int64, action string) (bool, error) {
	return false, nil
}
func (m *mockAccessService) CanViewCourse(ctx context.Context, course *domain.Course, userID int64) (bool, error) {
	return false, nil
}
func (m *mockAccessService) CanEditCourse(ctx context.Context, course *domain.Course, userID int64) (bool, error) {
	if m.err != nil {
		return false, m.err
	}
	return m.canEdit, nil
}
func (m *mockAccessService) CanDeleteCourse(ctx context.Context, course *domain.Course, userID int64) (bool, error) {
	return false, nil
}
func (m *mockAccessService) CanPublishCourse(ctx context.Context, course *domain.Course, userID int64) (bool, error) {
	return false, nil
}
func (m *mockAccessService) CanSelfEnrollCourse(ctx context.Context, course *domain.Course, userID int64) (bool, error) {
	return false, nil
}
func (m *mockAccessService) CanManageCourseUsers(ctx context.Context, course *domain.Course, userID int64) (bool, error) {
	return false, nil
}

func TestListTeacherPendingHomeworks_AllCourses_Success(t *testing.T) {
	repo := &mockAnalyticsRepo{
		items: []domain.PendingHomeworkItem{
			{
				AttemptID:       1,
				AnswerID:        10,
				StudentID:       5,
				StudentName:     "Test Student",
				CourseID:        2,
				CourseTitle:     "Russian",
				LessonID:        3,
				LessonTitle:     "Lesson 1",
				QuestionID:      4,
				QuestionText:    "Write essay",
				StudentAnswer:   "My answer",
				MaxPoints:       20,
				SubmittedAt:     time.Now(),
			},
		},
		total: 1,
		summary: []domain.CoursePendingSummaryItem{
			{CourseID: 2, CourseTitle: "Russian", PendingCount: 1},
			{CourseID: 3, CourseTitle: "Math", PendingCount: 0},
		},
	}
	courseRepo := &mockCourseRepo{}
	accessSvc := &mockAccessService{canEdit: true}

	svc := analyticsService.NewAnalyticsService(repo, courseRepo, accessSvc, nil)

	res, err := svc.ListTeacherPendingHomeworks(context.Background(), 1, 0, 1, 20)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res == nil {
		t.Fatal("expected result, got nil")
	}

	if res.Total != 1 {
		t.Errorf("expected total 1, got %d", res.Total)
	}

	if len(res.Items) != 1 {
		t.Errorf("expected 1 item, got %d", len(res.Items))
	}

	if len(res.CoursesSummary) != 2 {
		t.Errorf("expected 2 summary items, got %d", len(res.CoursesSummary))
	}
}

func TestListTeacherPendingHomeworks_SpecificCourse_AccessGranted(t *testing.T) {
	repo := &mockAnalyticsRepo{
		items: []domain.PendingHomeworkItem{
			{
				AttemptID:   1,
				CourseID:    2,
				SubmittedAt: time.Now(),
			},
		},
		total: 1,
		summary: []domain.CoursePendingSummaryItem{
			{CourseID: 2, CourseTitle: "Russian", PendingCount: 1},
		},
	}
	courseRepo := &mockCourseRepo{
		course: &domain.Course{Id: 2, Title: "Russian", CreatedBy: 1},
	}
	accessSvc := &mockAccessService{canEdit: true}

	svc := analyticsService.NewAnalyticsService(repo, courseRepo, accessSvc, nil)

	res, err := svc.ListTeacherPendingHomeworks(context.Background(), 1, 2, 1, 10)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.Total != 1 {
		t.Errorf("expected total 1, got %d", res.Total)
	}
}

func TestListTeacherPendingHomeworks_SpecificCourse_AccessDenied_Forbidden(t *testing.T) {
	repo := &mockAnalyticsRepo{}
	courseRepo := &mockCourseRepo{
		course: &domain.Course{Id: 99, Title: "Foreign Course", CreatedBy: 42},
	}
	accessSvc := &mockAccessService{canEdit: false}

	svc := analyticsService.NewAnalyticsService(repo, courseRepo, accessSvc, nil)

	_, err := svc.ListTeacherPendingHomeworks(context.Background(), 1, 99, 1, 10)
	if err == nil {
		t.Fatal("expected error for forbidden course access, got nil")
	}

	if !errors.Is(err, errorsAPP.ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestListTeacherPendingHomeworks_EmptyResults_NonNilSlices(t *testing.T) {
	repo := &mockAnalyticsRepo{
		items:   nil,
		total:   0,
		summary: nil,
	}
	courseRepo := &mockCourseRepo{}
	accessSvc := &mockAccessService{canEdit: true}

	svc := analyticsService.NewAnalyticsService(repo, courseRepo, accessSvc, nil)

	res, err := svc.ListTeacherPendingHomeworks(context.Background(), 1, 0, 1, 20)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.Items == nil {
		t.Error("expected non-nil Items slice")
	}
	if res.CoursesSummary == nil {
		t.Error("expected non-nil CoursesSummary slice")
	}
}
