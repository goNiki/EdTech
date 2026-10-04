package enrollment_test

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"

	"edtech/internal/domain"
	"edtech/internal/repository"
	enrollmentService "edtech/internal/service/enrollment"
	errorsAPP "edtech/pkg/errors"

	"github.com/jackc/pgx/v5"
)

type raceMockTxManager struct{}

func (m *raceMockTxManager) WithTX(ctx context.Context, opts pgx.TxOptions, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

type raceMockCourseRepo struct {
	repository.CourseRepository
	incrementCount atomic.Int64
}

func (m *raceMockCourseRepo) GetCourseByID(ctx context.Context, id int64) (*domain.Course, error) {
	return &domain.Course{
		Id:         id,
		Status:     domain.StatusPublished,
		Visibility: domain.VisibilityPublic,
		CreatedBy:  999,
	}, nil
}

func (m *raceMockCourseRepo) IncrementEnrolledCount(ctx context.Context, courseID int64) error {
	m.incrementCount.Add(1)
	return nil
}

type raceMockEnrollmentRepo struct {
	repository.EnrolledRepository
	mu       sync.Mutex
	enrolled map[string]bool
}

func newRaceMockEnrollmentRepo() *raceMockEnrollmentRepo {
	return &raceMockEnrollmentRepo{
		enrolled: make(map[string]bool),
	}
}

func (m *raceMockEnrollmentRepo) EnrollUserToCourse(ctx context.Context, enroll domain.EnrolledInCourse) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	key := enroll.Role + ":" + string(rune(enroll.UserID)) + ":" + string(rune(enroll.CourseID))
	if m.enrolled[key] {
		return errorsAPP.ErrUserAlreadyEnrolled
	}
	m.enrolled[key] = true
	return nil
}

func TestSelfEnrollCourse_ConcurrentRaceCondition(t *testing.T) {
	courseRepo := &raceMockCourseRepo{}
	enrolledRepo := newRaceMockEnrollmentRepo()
	txMgr := &raceMockTxManager{}

	svc := enrollmentService.NewEnrolmentService(courseRepo, enrolledRepo, nil, nil, txMgr)

	const concurrency = 20
	var wg sync.WaitGroup
	var successCount atomic.Int64
	var alreadyEnrolledCount atomic.Int64
	var unexpectedErrCount atomic.Int64

	startSignal := make(chan struct{})

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-startSignal

			err := svc.SelfEnrollCourse(context.Background(), domain.SelfEnrollRequest{
				UserID:   42,
				CourseID: 100,
			})

			if err == nil {
				successCount.Add(1)
			} else if errors.Is(err, errorsAPP.ErrUserAlreadyEnrolled) {
				alreadyEnrolledCount.Add(1)
			} else {
				unexpectedErrCount.Add(1)
			}
		}()
	}

	// Trigger all goroutines simultaneously
	close(startSignal)
	wg.Wait()

	if unexpectedErrCount.Load() > 0 {
		t.Fatalf("encountered %d unexpected errors", unexpectedErrCount.Load())
	}

	if success := successCount.Load(); success != 1 {
		t.Fatalf("expected exactly 1 successful enrollment, got %d", success)
	}

	if already := alreadyEnrolledCount.Load(); already != concurrency-1 {
		t.Fatalf("expected %d ErrUserAlreadyEnrolled, got %d", concurrency-1, already)
	}

	if count := courseRepo.incrementCount.Load(); count != 1 {
		t.Fatalf("expected IncrementEnrolledCount to be called exactly once, got %d", count)
	}
}
