package quiz_test

import (
	"context"
	"errors"
	"testing"

	"edtech/internal/domain"
	"edtech/internal/repository"
	"edtech/internal/service"
	"edtech/internal/service/quiz"
	errorsAPP "edtech/pkg/errors"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockQuizRepoForGradingValidation struct {
	repository.QuizRepository
	attempt          *domain.QuizAttempt
	quiz             *domain.Quiz
	maxPoints        int
	updateCalled     bool
	updatedPoints    int
	ungradedCount    int
	sumPoints        int
}

func (m *mockQuizRepoForGradingValidation) GetAttemptForUpdate(ctx context.Context, attemptID int64) (*domain.QuizAttempt, error) {
	if m.attempt != nil {
		return m.attempt, nil
	}
	return nil, errorsAPP.ErrAttemptNotFound
}

func (m *mockQuizRepoForGradingValidation) GetQuizByID(ctx context.Context, quizID int64) (*domain.Quiz, error) {
	if m.quiz != nil {
		return m.quiz, nil
	}
	return nil, errorsAPP.ErrQuizNotFound
}

func (m *mockQuizRepoForGradingValidation) GetQuizTotalPoints(ctx context.Context, quizID int64) (int, error) {
	return m.maxPoints, nil
}

func (m *mockQuizRepoForGradingValidation) UpdateAttemptAnswer(ctx context.Context, answerID int64, points int, feedback *string, isCorrect bool, gradedBy ...int64) error {
	m.updateCalled = true
	m.updatedPoints = points
	return nil
}

func (m *mockQuizRepoForGradingValidation) CountUngradedAnswers(ctx context.Context, attemptID int64) (int, error) {
	return m.ungradedCount, nil
}

func (m *mockQuizRepoForGradingValidation) SumAttemptPoints(ctx context.Context, attemptID int64) (int, error) {
	return m.sumPoints, nil
}

func (m *mockQuizRepoForGradingValidation) UpdateAttempt(ctx context.Context, attempt *domain.QuizAttempt) error {
	return nil
}

type mockLessonRepoForGradingValidation struct {
	repository.LessonRepository
	lesson *domain.Lesson
}

func (m *mockLessonRepoForGradingValidation) GetLessonByID(ctx context.Context, id int64) (*domain.Lesson, error) {
	if m.lesson != nil {
		return m.lesson, nil
	}
	return nil, errorsAPP.ErrLessonNotFound
}

type mockCourseRepoForGradingValidation struct {
	repository.CourseRepository
	course *domain.Course
}

func (m *mockCourseRepoForGradingValidation) GetCourseByID(ctx context.Context, id int64) (*domain.Course, error) {
	if m.course != nil {
		return m.course, nil
	}
	return nil, errorsAPP.ErrCourseNotFound
}

type mockAccessSvcForGradingValidation struct {
	service.AccessService
	canEdit bool
}

func (m *mockAccessSvcForGradingValidation) CanEditCourse(ctx context.Context, course *domain.Course, userID int64) (bool, error) {
	return m.canEdit, nil
}

type mockTxManagerForGradingValidation struct{}

func (m *mockTxManagerForGradingValidation) WithTX(ctx context.Context, opts pgx.TxOptions, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func setupGradingService(maxPoints int) (service.QuizServices, *mockQuizRepoForGradingValidation) {
	attempt := &domain.QuizAttempt{
		ID:     18,
		QuizID: 10,
		UserID: 5,
	}
	quizObj := &domain.Quiz{
		ID:          10,
		LessonID:    1,
		PassingScor: 70,
	}
	lesson := &domain.Lesson{
		ID:       1,
		CourseID: 2,
	}
	course := &domain.Course{
		Id:        2,
		CreatedBy: 1,
	}

	quizRepo := &mockQuizRepoForGradingValidation{
		attempt:       attempt,
		quiz:          quizObj,
		maxPoints:     maxPoints,
		ungradedCount: 1,
	}
	lessonRepo := &mockLessonRepoForGradingValidation{lesson: lesson}
	courseRepo := &mockCourseRepoForGradingValidation{course: course}
	accessSvc := &mockAccessSvcForGradingValidation{canEdit: true}
	txMgr := &mockTxManagerForGradingValidation{}

	svc := quiz.NewQuizService(
		quizRepo,
		courseRepo,
		lessonRepo,
		accessSvc,
		nil,
		txMgr,
		nil,
	)

	return svc, quizRepo
}

func TestGradeAttemptAnswer_NegativePoints_Rejected(t *testing.T) {
	svc, quizRepo := setupGradingService(25)

	attempt, err := svc.GradeAttemptAnswer(context.Background(), 1, 18, 16, -5, nil)

	require.Error(t, err)
	assert.Nil(t, attempt)
	assert.True(t, errors.Is(err, errorsAPP.ErrInvalidGradePoints), "expected ErrInvalidGradePoints, got %v", err)
	assert.False(t, quizRepo.updateCalled, "UpdateAttemptAnswer should not be called")
}

func TestGradeAttemptAnswer_PointsExceedMax_Rejected(t *testing.T) {
	svc, quizRepo := setupGradingService(25)

	attempt, err := svc.GradeAttemptAnswer(context.Background(), 1, 18, 16, 999, nil)

	require.Error(t, err)
	assert.Nil(t, attempt)
	assert.True(t, errors.Is(err, errorsAPP.ErrInvalidGradePoints), "expected ErrInvalidGradePoints, got %v", err)
	assert.Contains(t, err.Error(), "25")
	assert.False(t, quizRepo.updateCalled, "UpdateAttemptAnswer should not be called")
}

func TestGradeAttemptAnswer_ValidPoints_Success(t *testing.T) {
	svc, quizRepo := setupGradingService(25)

	feedback := "Хороший ответ"
	attempt, err := svc.GradeAttemptAnswer(context.Background(), 1, 18, 16, 20, &feedback)

	require.NoError(t, err)
	assert.NotNil(t, attempt)
	assert.True(t, quizRepo.updateCalled, "UpdateAttemptAnswer should be called")
	assert.Equal(t, 20, quizRepo.updatedPoints)
}
