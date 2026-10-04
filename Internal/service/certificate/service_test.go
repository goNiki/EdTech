package certificate

import (
	"context"
	"errors"
	"testing"
	"time"

	"edtech/internal/domain"
	errorsAPP "edtech/pkg/errors"
)

type mockCertRepo struct {
	certsByCode   map[string]*domain.Certificate
	certsByUser   map[string]*domain.Certificate
	createCallCnt int
}

func newMockCertRepo() *mockCertRepo {
	return &mockCertRepo{
		certsByCode: make(map[string]*domain.Certificate),
		certsByUser: make(map[string]*domain.Certificate),
	}
}

func (m *mockCertRepo) CreateCertificate(ctx context.Context, cert *domain.Certificate) (*domain.Certificate, error) {
	m.createCallCnt++
	copyCert := *cert
	copyCert.ID = 100
	copyCert.IssuedAt = time.Now()
	key := fmtUserCourse(cert.UserID, cert.CourseID)
	m.certsByUser[key] = &copyCert
	m.certsByCode[copyCert.CertificateCode] = &copyCert
	return &copyCert, nil
}

func (m *mockCertRepo) GetCertificateByCode(ctx context.Context, code string) (*domain.Certificate, error) {
	cert, ok := m.certsByCode[code]
	if !ok {
		return nil, errorsAPP.ErrCertificateNotFound
	}
	return cert, nil
}

func (m *mockCertRepo) GetCertificateByUserAndCourse(ctx context.Context, userID, courseID int64) (*domain.Certificate, error) {
	key := fmtUserCourse(userID, courseID)
	cert, ok := m.certsByUser[key]
	if !ok {
		return nil, errorsAPP.ErrCertificateNotFound
	}
	return cert, nil
}

func fmtUserCourse(u, c int64) string {
	return string(rune(u)) + ":" + string(rune(c))
}

type mockProgressRepo struct {
	progress *domain.CourseProgress
	err      error
}

func (m *mockProgressRepo) CreateLessonProgress(ctx context.Context, progress *domain.LessonProgress) error {
	return nil
}
func (m *mockProgressRepo) GetLessonProgress(ctx context.Context, userID, lessonID int64) (*domain.LessonProgress, error) {
	return nil, nil
}
func (m *mockProgressRepo) UpdateLessonProgressTime(ctx context.Context, userID, lessonID int64, additionalTime int, lastPos int) error {
	return nil
}
func (m *mockProgressRepo) UpdateLessonProgressStatus(ctx context.Context, userID, lessonID int64, status domain.ProgressStatus) error {
	return nil
}
func (m *mockProgressRepo) UpdateLessonProgressScore(ctx context.Context, userID, lessonID int64, score int) error {
	return nil
}
func (m *mockProgressRepo) CreateCourseProgress(ctx context.Context, progress *domain.CourseProgress) error {
	return nil
}
func (m *mockProgressRepo) UpsertCourseProgress(ctx context.Context, userID, courseID int64, completedLessons, totalLessons int, percentage float64) error {
	return nil
}
func (m *mockProgressRepo) UpsertCourseProgressWithScore(ctx context.Context, userID, courseID int64, completedLessons, totalLessons int, percentage float64, averageScore float64) error {
	return nil
}
func (m *mockProgressRepo) GetCourseProgress(ctx context.Context, userID, courseID int64) (*domain.CourseProgress, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.progress, nil
}
func (m *mockProgressRepo) GetAllLessonProgressByCourse(ctx context.Context, userID, courseID int64) ([]domain.LessonProgress, error) {
	return nil, nil
}

type mockUserRepo struct {
	user *domain.User
}

func (m *mockUserRepo) CreateUser(ctx context.Context, user domain.CreateUser) (*domain.User, error) {
	return nil, nil
}
func (m *mockUserRepo) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return nil, nil
}
func (m *mockUserRepo) GetUserByUserName(ctx context.Context, username string) (*domain.User, error) {
	return nil, nil
}
func (m *mockUserRepo) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	return m.user, nil
}
func (m *mockUserRepo) ExistingByEmail(ctx context.Context, email string) (bool, error) {
	return false, nil
}
func (m *mockUserRepo) ExistingByUsernName(ctx context.Context, username string) (bool, error) {
	return false, nil
}
func (m *mockUserRepo) UpdateLastLogin(ctx context.Context, userID int64, now time.Time) error {
	return nil
}
func (m *mockUserRepo) UpdateProfile(ctx context.Context, user *domain.User) error {
	return nil
}
func (m *mockUserRepo) UpdatePassword(ctx context.Context, userID int64, passHash string) error {
	return nil
}
func (m *mockUserRepo) SetEmailVerified(ctx context.Context, userID int64, verified bool) error {
	return nil
}
func (m *mockUserRepo) UpdateRole(ctx context.Context, userID int64, role domain.Role) error {
	return nil
}
func (m *mockUserRepo) SetBannedStatus(ctx context.Context, userID int64, isBanned bool) error {
	return nil
}
func (m *mockUserRepo) ListUsers(ctx context.Context, filter domain.UserFilter, pagination domain.Pagination) ([]domain.User, int64, error) {
	return nil, 0, nil
}

type mockCourseRepo struct {
	course *domain.Course
}

func (m *mockCourseRepo) CreateCourse(ctx context.Context, course *domain.Course) (*domain.Course, error) {
	return nil, nil
}
func (m *mockCourseRepo) GetCourseByID(ctx context.Context, id int64) (*domain.Course, error) {
	return m.course, nil
}
func (m *mockCourseRepo) GetCourseBySlug(ctx context.Context, slug string) (*domain.Course, error) {
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

func TestGetOrIssueCertificate_Success(t *testing.T) {
	certRepo := newMockCertRepo()
	avgScore := 95.5
	progressRepo := &mockProgressRepo{
		progress: &domain.CourseProgress{
			Percent:       100,
			TotalLessons:  10,
			CompletedLess: 10,
			AverageScore:  &avgScore,
		},
	}
	firstName := "Ivan"
	lastName := "Ivanov"
	userRepo := &mockUserRepo{
		user: &domain.User{
			ID:        1,
			Username:  "ivan",
			FirstName: &firstName,
			LastName:  &lastName,
		},
	}
	courseRepo := &mockCourseRepo{
		course: &domain.Course{
			Id:    10,
			Title: "Golang Professional",
		},
	}

	svc := NewCertificateService(certRepo, courseRepo, userRepo, progressRepo)

	cert, err := svc.GetOrIssueCertificate(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cert == nil {
		t.Fatal("expected certificate, got nil")
	}
	if cert.StudentName != "Ivan Ivanov" {
		t.Errorf("expected student name 'Ivan Ivanov', got '%s'", cert.StudentName)
	}
	if cert.CourseTitle != "Golang Professional" {
		t.Errorf("expected course title 'Golang Professional', got '%s'", cert.CourseTitle)
	}
	if cert.FinalScore != 95.5 {
		t.Errorf("expected score 95.5, got %v", cert.FinalScore)
	}
	if cert.CertificateCode == "" {
		t.Error("expected non-empty certificate code")
	}
	if certRepo.createCallCnt != 1 {
		t.Errorf("expected 1 create call, got %d", certRepo.createCallCnt)
	}

	// Idempotency check: second call returns the same certificate
	secondCert, err := svc.GetOrIssueCertificate(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("unexpected error on second call: %v", err)
	}
	if secondCert.CertificateCode != cert.CertificateCode {
		t.Errorf("expected same code on second call, got %s vs %s", secondCert.CertificateCode, cert.CertificateCode)
	}
	if certRepo.createCallCnt != 1 {
		t.Errorf("expected still 1 create call due to idempotency, got %d", certRepo.createCallCnt)
	}
}

func TestGetOrIssueCertificate_CourseNotCompleted(t *testing.T) {
	certRepo := newMockCertRepo()
	progressRepo := &mockProgressRepo{
		progress: &domain.CourseProgress{
			Percent:       80,
			TotalLessons:  10,
			CompletedLess: 8,
		},
	}
	svc := NewCertificateService(certRepo, &mockCourseRepo{}, &mockUserRepo{}, progressRepo)

	_, err := svc.GetOrIssueCertificate(context.Background(), 1, 10)
	if !errors.Is(err, errorsAPP.ErrCourseNotCompleted) {
		t.Fatalf("expected ErrCourseNotCompleted, got: %v", err)
	}
}

func TestVerifyCertificate(t *testing.T) {
	certRepo := newMockCertRepo()
	code := "EDL-2026-TESTCODE"
	certRepo.certsByCode[code] = &domain.Certificate{
		CertificateCode: code,
		StudentName:     "Test Student",
		CourseTitle:     "Web Dev",
		FinalScore:      100,
	}

	svc := NewCertificateService(certRepo, &mockCourseRepo{}, &mockUserRepo{}, &mockProgressRepo{})

	// Valid verify
	cert, err := svc.VerifyCertificate(context.Background(), code)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cert.StudentName != "Test Student" {
		t.Errorf("expected student name 'Test Student', got '%s'", cert.StudentName)
	}

	// Unknown verify
	_, err = svc.VerifyCertificate(context.Background(), "EDL-9999-NOTFOUND")
	if !errors.Is(err, errorsAPP.ErrCertificateNotFound) {
		t.Fatalf("expected ErrCertificateNotFound, got: %v", err)
	}
}
