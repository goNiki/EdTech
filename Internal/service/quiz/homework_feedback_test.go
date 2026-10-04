package quiz_test

import (
	"context"
	"testing"
	"time"

	"edtech/internal/domain"
	"edtech/internal/repository"
	"edtech/internal/service"
	"edtech/internal/service/quiz"
	errorsAPP "edtech/pkg/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockQuizRepoForFeedback struct {
	repository.QuizRepository
	feedbackResp *domain.StudentHomeworkFeedback
	err          error
}

func (m *mockQuizRepoForFeedback) GetStudentHomeworkFeedback(ctx context.Context, userID, lessonID int64) (*domain.StudentHomeworkFeedback, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.feedbackResp, nil
}

type mockLessonRepoForFeedback struct {
	repository.LessonRepository
	lesson *domain.Lesson
	err    error
}

func (m *mockLessonRepoForFeedback) GetLessonByID(ctx context.Context, id int64) (*domain.Lesson, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.lesson != nil {
		return m.lesson, nil
	}
	return nil, errorsAPP.ErrLessonNotFound
}

type mockCourseRepoForFeedback struct {
	repository.CourseRepository
	course *domain.Course
	err    error
}

func (m *mockCourseRepoForFeedback) GetCourseByID(ctx context.Context, id int64) (*domain.Course, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.course != nil {
		return m.course, nil
	}
	return nil, errorsAPP.ErrCourseNotFound
}

type mockAccessServiceForFeedback struct {
	service.AccessService
	canView bool
	err     error
}

func (m *mockAccessServiceForFeedback) CanViewCourse(ctx context.Context, course *domain.Course, userID int64) (bool, error) {
	if m.err != nil {
		return false, m.err
	}
	return m.canView, nil
}

func TestGetStudentHomeworkFeedback_Graded(t *testing.T) {
	ctx := context.Background()

	lessonID := int64(10)
	courseID := int64(1)
	studentID := int64(5)
	teacherID := int64(2)
	avatar := "/static/avatar.png"
	feedbackText := "Отличная работа!"
	now := time.Now()
	attemptID := int64(100)
	isCorrect := true

	qRepo := &mockQuizRepoForFeedback{
		feedbackResp: &domain.StudentHomeworkFeedback{
			HasSubmission: true,
			Status:        "graded",
			AttemptID:     &attemptID,
			SubmittedAt:   &now,
			GradedAt:      &now,
			Teacher: &domain.HomeworkTeacherInfo{
				ID:        teacherID,
				Name:      "Преподаватель Тест",
				AvatarURL: &avatar,
			},
			Answers: []domain.HomeworkAnswerDetail{
				{
					AnswerID:      1,
					QuestionText:  "Напишите эссе",
					StudentAnswer: "Мое эссе...",
					Points:        25,
					MaxPoints:     25,
					IsCorrect:     &isCorrect,
					Feedback:      &feedbackText,
				},
			},
		},
	}

	lRepo := &mockLessonRepoForFeedback{
		lesson: &domain.Lesson{ID: lessonID, CourseID: courseID},
	}
	cRepo := &mockCourseRepoForFeedback{
		course: &domain.Course{Id: courseID},
	}
	accessSvc := &mockAccessServiceForFeedback{canView: true}

	svc := quiz.NewQuizService(qRepo, cRepo, lRepo, accessSvc, nil, nil, nil)

	res, err := svc.GetStudentHomeworkFeedback(ctx, studentID, lessonID)
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.True(t, res.HasSubmission)
	assert.Equal(t, "graded", res.Status)
	assert.Equal(t, attemptID, *res.AttemptID)
	assert.NotNil(t, res.Teacher)
	assert.Equal(t, teacherID, res.Teacher.ID)
	assert.Equal(t, "Преподаватель Тест", res.Teacher.Name)
	assert.Len(t, res.Answers, 1)
	assert.Equal(t, 25, res.Answers[0].Points)
	assert.Equal(t, "Отличная работа!", *res.Answers[0].Feedback)
}

func TestGetStudentHomeworkFeedback_Pending(t *testing.T) {
	ctx := context.Background()

	lessonID := int64(10)
	courseID := int64(1)
	studentID := int64(5)
	attemptID := int64(101)
	now := time.Now()

	qRepo := &mockQuizRepoForFeedback{
		feedbackResp: &domain.StudentHomeworkFeedback{
			HasSubmission: true,
			Status:        "pending",
			AttemptID:     &attemptID,
			SubmittedAt:   &now,
			GradedAt:      nil,
			Teacher:       nil,
			Answers: []domain.HomeworkAnswerDetail{
				{
					AnswerID:      2,
					QuestionText:  "Напишите эссе",
					StudentAnswer: "Черновик ответа...",
					Points:        0,
					MaxPoints:     25,
					IsCorrect:     nil,
					Feedback:      nil,
				},
			},
		},
	}

	lRepo := &mockLessonRepoForFeedback{
		lesson: &domain.Lesson{ID: lessonID, CourseID: courseID},
	}
	cRepo := &mockCourseRepoForFeedback{
		course: &domain.Course{Id: courseID},
	}
	accessSvc := &mockAccessServiceForFeedback{canView: true}

	svc := quiz.NewQuizService(qRepo, cRepo, lRepo, accessSvc, nil, nil, nil)

	res, err := svc.GetStudentHomeworkFeedback(ctx, studentID, lessonID)
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.True(t, res.HasSubmission)
	assert.Equal(t, "pending", res.Status)
	assert.Nil(t, res.GradedAt)
	assert.Nil(t, res.Teacher)
	assert.Len(t, res.Answers, 1)
	assert.Nil(t, res.Answers[0].IsCorrect)
}

func TestGetStudentHomeworkFeedback_NotSubmitted(t *testing.T) {
	ctx := context.Background()

	lessonID := int64(10)
	courseID := int64(1)
	studentID := int64(5)

	qRepo := &mockQuizRepoForFeedback{
		feedbackResp: &domain.StudentHomeworkFeedback{
			HasSubmission: false,
			Status:        "not_submitted",
			Answers:       []domain.HomeworkAnswerDetail{},
		},
	}

	lRepo := &mockLessonRepoForFeedback{
		lesson: &domain.Lesson{ID: lessonID, CourseID: courseID},
	}
	cRepo := &mockCourseRepoForFeedback{
		course: &domain.Course{Id: courseID},
	}
	accessSvc := &mockAccessServiceForFeedback{canView: true}

	svc := quiz.NewQuizService(qRepo, cRepo, lRepo, accessSvc, nil, nil, nil)

	res, err := svc.GetStudentHomeworkFeedback(ctx, studentID, lessonID)
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.False(t, res.HasSubmission)
	assert.Equal(t, "not_submitted", res.Status)
	assert.Empty(t, res.Answers)
}

func TestGetStudentHomeworkFeedback_Unauthorized(t *testing.T) {
	ctx := context.Background()

	svc := quiz.NewQuizService(nil, nil, nil, nil, nil, nil, nil)

	_, err := svc.GetStudentHomeworkFeedback(ctx, 0, 10)
	require.Error(t, err)
	assert.ErrorIs(t, err, errorsAPP.ErrUnauthorized)
}

func TestGetStudentHomeworkFeedback_InvalidLessonID(t *testing.T) {
	ctx := context.Background()

	svc := quiz.NewQuizService(nil, nil, nil, nil, nil, nil, nil)

	_, err := svc.GetStudentHomeworkFeedback(ctx, 5, 0)
	require.Error(t, err)
	assert.ErrorIs(t, err, errorsAPP.ErrInvalidURLParam)
}

func TestGetStudentHomeworkFeedback_LessonNotFound(t *testing.T) {
	ctx := context.Background()

	lRepo := &mockLessonRepoForFeedback{err: errorsAPP.ErrLessonNotFound}
	svc := quiz.NewQuizService(nil, nil, lRepo, nil, nil, nil, nil)

	_, err := svc.GetStudentHomeworkFeedback(ctx, 5, 999)
	require.Error(t, err)
	assert.ErrorIs(t, err, errorsAPP.ErrLessonNotFound)
}

func TestGetStudentHomeworkFeedback_Forbidden(t *testing.T) {
	ctx := context.Background()

	lessonID := int64(10)
	courseID := int64(1)

	lRepo := &mockLessonRepoForFeedback{
		lesson: &domain.Lesson{ID: lessonID, CourseID: courseID},
	}
	cRepo := &mockCourseRepoForFeedback{
		course: &domain.Course{Id: courseID},
	}
	accessSvc := &mockAccessServiceForFeedback{canView: false}

	svc := quiz.NewQuizService(nil, cRepo, lRepo, accessSvc, nil, nil, nil)

	_, err := svc.GetStudentHomeworkFeedback(ctx, 5, lessonID)
	require.Error(t, err)
	assert.ErrorIs(t, err, errorsAPP.ErrForbidden)
}
