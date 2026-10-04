package review_test

import (
	"context"
	"testing"

	"edtech/internal/domain"
	"edtech/internal/infrastructure/db"
	"edtech/internal/service/review"
	errorsAPP "edtech/pkg/errors"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockReviewRepo struct {
	upsertFunc    func(ctx context.Context, q db.QueryExecutor, r *domain.Review) (*domain.Review, error)
	deleteFunc    func(ctx context.Context, q db.QueryExecutor, courseID, userID int64) error
	getReviewFunc func(ctx context.Context, q db.QueryExecutor, courseID, userID int64) (*domain.Review, error)
	listFunc      func(ctx context.Context, q db.QueryExecutor, courseID int64, limit, offset int) ([]domain.Review, int64, error)
	summaryFunc   func(ctx context.Context, q db.QueryExecutor, courseID int64) (float64, int, error)
}

func (m *mockReviewRepo) UpsertReview(ctx context.Context, q db.QueryExecutor, r *domain.Review) (*domain.Review, error) {
	if m.upsertFunc != nil {
		return m.upsertFunc(ctx, q, r)
	}
	return r, nil
}
func (m *mockReviewRepo) DeleteReview(ctx context.Context, q db.QueryExecutor, courseID, userID int64) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, q, courseID, userID)
	}
	return nil
}
func (m *mockReviewRepo) GetReviewByUserAndCourse(ctx context.Context, q db.QueryExecutor, courseID, userID int64) (*domain.Review, error) {
	if m.getReviewFunc != nil {
		return m.getReviewFunc(ctx, q, courseID, userID)
	}
	return nil, nil
}
func (m *mockReviewRepo) ListReviewsByCourse(ctx context.Context, q db.QueryExecutor, courseID int64, limit, offset int) ([]domain.Review, int64, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, q, courseID, limit, offset)
	}
	return nil, 0, nil
}
func (m *mockReviewRepo) GetCourseRatingSummary(ctx context.Context, q db.QueryExecutor, courseID int64) (float64, int, error) {
	if m.summaryFunc != nil {
		return m.summaryFunc(ctx, q, courseID)
	}
	return 4.5, 10, nil
}

type mockCourseRepo struct {
	updateStatsFunc func(ctx context.Context, q db.QueryExecutor, courseID int64, rating float64, count int) error
}

func (m *mockCourseRepo) GetCourseBySlug(ctx context.Context, q db.QueryExecutor, slug string) (*domain.Course, error) {
	return nil, nil
}
func (m *mockCourseRepo) CreateCourse(ctx context.Context, q db.QueryExecutor, course *domain.Course) (*domain.Course, error) {
	return nil, nil
}
func (m *mockCourseRepo) GetCourseByID(ctx context.Context, q db.QueryExecutor, id int64) (*domain.Course, error) {
	return nil, nil
}
func (m *mockCourseRepo) ArchiveCourse(ctx context.Context, q db.QueryExecutor, course *domain.Course) error {
	return nil
}
func (m *mockCourseRepo) PublishCourse(ctx context.Context, q db.QueryExecutor, course *domain.Course) error {
	return nil
}
func (m *mockCourseRepo) UpdateCourseStatus(ctx context.Context, q db.QueryExecutor, courseID int64, status string) error {
	return nil
}
func (m *mockCourseRepo) UpdateCourse(ctx context.Context, q db.QueryExecutor, course *domain.Course) error {
	return nil
}
func (m *mockCourseRepo) ListPublicCourses(ctx context.Context, q db.QueryExecutor, pagination domain.Pagination, filter domain.CourseFilter) ([]domain.Course, error) {
	return nil, nil
}
func (m *mockCourseRepo) ListEnrolledCourses(ctx context.Context, q db.QueryExecutor, input *domain.InputListMyCourse) ([]domain.Course, error) {
	return nil, nil
}
func (m *mockCourseRepo) CountEnrolledCourses(ctx context.Context, q db.QueryExecutor, input *domain.InputListMyCourse) (int64, error) {
	return 0, nil
}
func (m *mockCourseRepo) CountCourses(ctx context.Context, q db.QueryExecutor, filter domain.CourseFilter) (int64, error) {
	return 0, nil
}
func (m *mockCourseRepo) DeleteCourse(ctx context.Context, q db.QueryExecutor, courseID int64) error {
	return nil
}
func (m *mockCourseRepo) ExistingBySlug(ctx context.Context, q db.QueryExecutor, slug string) (bool, error) {
	return false, nil
}
func (m *mockCourseRepo) ExistsBySlug(ctx context.Context, q db.QueryExecutor, slug string) (bool, error) {
	return false, nil
}
func (m *mockCourseRepo) IncrementEnrolledCount(ctx context.Context, q db.QueryExecutor, courseID int64) error {
	return nil
}
func (m *mockCourseRepo) DecrementEnrolledCount(ctx context.Context, q db.QueryExecutor, courseID int64) error {
	return nil
}
func (m *mockCourseRepo) UpdateCourseRatingStats(ctx context.Context, q db.QueryExecutor, courseID int64, rating float64, count int) error {
	if m.updateStatsFunc != nil {
		return m.updateStatsFunc(ctx, q, courseID, rating, count)
	}
	return nil
}

type mockEnrolledRepo struct {
	userExistFunc func(ctx context.Context, q db.QueryExecutor, userID, courseID int64) (bool, error)
}

func (m *mockEnrolledRepo) UserExistCourse(ctx context.Context, q db.QueryExecutor, userID, courseID int64) (bool, error) {
	if m.userExistFunc != nil {
		return m.userExistFunc(ctx, q, userID, courseID)
	}
	return true, nil
}
func (m *mockEnrolledRepo) EnrollUserToCourse(ctx context.Context, q db.QueryExecutor, enroll domain.EnrolledInCourse) error {
	return nil
}
func (m *mockEnrolledRepo) GetRoleUserInCourse(ctx context.Context, q db.QueryExecutor, userID, courseID int64) (string, error) {
	return "student", nil
}
func (m *mockEnrolledRepo) UnenrollUser(ctx context.Context, q db.QueryExecutor, userID, courseID int64) error {
	return nil
}
func (m *mockEnrolledRepo) ListCourseStudents(ctx context.Context, q db.QueryExecutor, courseID int64, limit, offset int64) ([]domain.User, int, error) {
	return nil, 0, nil
}
func (m *mockEnrolledRepo) ListCourseStudentsWithProgress(ctx context.Context, q db.QueryExecutor, courseID int64, limit, offset int64) ([]domain.CourseStudentItem, int64, error) {
	return nil, 0, nil
}
func (m *mockEnrolledRepo) ChangeUserRole(ctx context.Context, q db.QueryExecutor, courseID int64, targetUserID int64, newRole string) error {
	return nil
}

type mockProgressRepo struct {
	getCourseProgressFunc func(ctx context.Context, q db.QueryExecutor, userID, courseID int64) (*domain.CourseProgress, error)
}

func (m *mockProgressRepo) CreateLessonProgress(ctx context.Context, q db.QueryExecutor, progress *domain.LessonProgress) error {
	return nil
}
func (m *mockProgressRepo) GetLessonProgress(ctx context.Context, q db.QueryExecutor, userID, lessonID int64) (*domain.LessonProgress, error) {
	return nil, nil
}
func (m *mockProgressRepo) UpdateLessonProgressTime(ctx context.Context, q db.QueryExecutor, userID, lessonID int64, additionalTime int, lastPos int) error {
	return nil
}
func (m *mockProgressRepo) UpdateLessonProgressStatus(ctx context.Context, q db.QueryExecutor, userID, lessonID int64, status domain.ProgressStatus) error {
	return nil
}
func (m *mockProgressRepo) UpdateLessonProgressScore(ctx context.Context, q db.QueryExecutor, userID, lessonID int64, score int) error {
	return nil
}
func (m *mockProgressRepo) CreateCourseProgress(ctx context.Context, q db.QueryExecutor, progress *domain.CourseProgress) error {
	return nil
}
func (m *mockProgressRepo) UpsertCourseProgress(ctx context.Context, q db.QueryExecutor, userID, courseID int64, completedLessons, totalLessons int, percentage float64) error {
	return nil
}
func (m *mockProgressRepo) UpsertCourseProgressWithScore(ctx context.Context, q db.QueryExecutor, userID, courseID int64, completedLessons, totalLessons int, percentage float64, averageScore float64) error {
	return nil
}
func (m *mockProgressRepo) GetCourseProgress(ctx context.Context, q db.QueryExecutor, userID, courseID int64) (*domain.CourseProgress, error) {
	if m.getCourseProgressFunc != nil {
		return m.getCourseProgressFunc(ctx, q, userID, courseID)
	}
	return &domain.CourseProgress{Percent: 50}, nil
}
func (m *mockProgressRepo) GetAllLessonProgressByCourse(ctx context.Context, q db.QueryExecutor, userID, courseID int64) ([]domain.LessonProgress, error) {
	return nil, nil
}

type mockTxManager struct{}

func (m *mockTxManager) WithTX(ctx context.Context, opts pgx.TxOptions, fn func(ctx context.Context, tx db.QueryExecutor) error) error {
	return fn(ctx, nil)
}

func TestAddOrUpdateReview_Success(t *testing.T) {
	reviewRepo := &mockReviewRepo{
		upsertFunc: func(ctx context.Context, q db.QueryExecutor, r *domain.Review) (*domain.Review, error) {
			r.ID = 101
			return r, nil
		},
		summaryFunc: func(ctx context.Context, q db.QueryExecutor, courseID int64) (float64, int, error) {
			return 4.8, 5, nil
		},
	}
	var updatedRating float64
	var updatedCount int
	courseRepo := &mockCourseRepo{
		updateStatsFunc: func(ctx context.Context, q db.QueryExecutor, courseID int64, rating float64, count int) error {
			updatedRating = rating
			updatedCount = count
			return nil
		},
	}
	enrolledRepo := &mockEnrolledRepo{
		userExistFunc: func(ctx context.Context, q db.QueryExecutor, userID, courseID int64) (bool, error) {
			return true, nil
		},
	}
	progressRepo := &mockProgressRepo{
		getCourseProgressFunc: func(ctx context.Context, q db.QueryExecutor, userID, courseID int64) (*domain.CourseProgress, error) {
			return &domain.CourseProgress{Percent: 35}, nil
		},
	}

	svc := review.NewReviewService(reviewRepo, courseRepo, enrolledRepo, progressRepo, &mockTxManager{}, nil)

	comment := "Great course!"
	rev, err := svc.AddOrUpdateReview(context.Background(), 1, 10, 5, &comment)
	require.NoError(t, err)
	require.NotNil(t, rev)
	assert.Equal(t, int64(101), rev.ID)
	assert.Equal(t, 5, rev.Rating)
	assert.Equal(t, 4.8, updatedRating)
	assert.Equal(t, 5, updatedCount)
}

func TestAddOrUpdateReview_ForbiddenIfNotEnrolled(t *testing.T) {
	enrolledRepo := &mockEnrolledRepo{
		userExistFunc: func(ctx context.Context, q db.QueryExecutor, userID, courseID int64) (bool, error) {
			return false, nil
		},
	}
	svc := review.NewReviewService(&mockReviewRepo{}, &mockCourseRepo{}, enrolledRepo, &mockProgressRepo{}, &mockTxManager{}, nil)

	_, err := svc.AddOrUpdateReview(context.Background(), 1, 10, 5, nil)
	require.Error(t, err)
	assert.ErrorIs(t, err, errorsAPP.ErrForbidden)
}

func TestAddOrUpdateReview_ForbiddenIfProgressBelowThreshold(t *testing.T) {
	enrolledRepo := &mockEnrolledRepo{
		userExistFunc: func(ctx context.Context, q db.QueryExecutor, userID, courseID int64) (bool, error) {
			return true, nil
		},
	}
	progressRepo := &mockProgressRepo{
		getCourseProgressFunc: func(ctx context.Context, q db.QueryExecutor, userID, courseID int64) (*domain.CourseProgress, error) {
			return &domain.CourseProgress{Percent: 29}, nil // Less than 30%
		},
	}
	svc := review.NewReviewService(&mockReviewRepo{}, &mockCourseRepo{}, enrolledRepo, progressRepo, &mockTxManager{}, nil)

	_, err := svc.AddOrUpdateReview(context.Background(), 1, 10, 5, nil)
	require.Error(t, err)
	assert.ErrorIs(t, err, errorsAPP.ErrForbidden)
}

func TestAddOrUpdateReview_InvalidRating(t *testing.T) {
	svc := review.NewReviewService(&mockReviewRepo{}, &mockCourseRepo{}, &mockEnrolledRepo{}, &mockProgressRepo{}, &mockTxManager{}, nil)

	_, err := svc.AddOrUpdateReview(context.Background(), 1, 10, 0, nil)
	require.Error(t, err)
	assert.ErrorIs(t, err, errorsAPP.ErrValidationFailed)

	_, err = svc.AddOrUpdateReview(context.Background(), 1, 10, 6, nil)
	require.Error(t, err)
	assert.ErrorIs(t, err, errorsAPP.ErrValidationFailed)
}
