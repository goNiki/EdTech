package progress_test

import (
	"context"
	"errors"
	"testing"

	"edtech/internal/domain"
	"edtech/internal/repository"
	progressService "edtech/internal/service/progress"
	errorsAPP "edtech/pkg/errors"

	"github.com/jackc/pgx/v5"
)

type mockProgressRepo struct {
	repository.ProgressRepository
	lessonProgress *domain.LessonProgress
	getErr         error
}

func (m *mockProgressRepo) GetLessonProgress(ctx context.Context, userID, lessonID int64) (*domain.LessonProgress, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	return m.lessonProgress, nil
}

type mockQuizRepo struct {
	repository.QuizRepository
	submissions []domain.LessonSubmissionDetail
	getErr      error
}

func (m *mockQuizRepo) GetLessonSubmissions(ctx context.Context, userID, lessonID int64) ([]domain.LessonSubmissionDetail, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	return m.submissions, nil
}

func TestGetLessonProgress_NotFound_ReturnsDefaultNotStarted(t *testing.T) {
	ctx := context.Background()
	pRepo := &mockProgressRepo{
		getErr: errorsAPP.ErrLessonProgressNotFound,
	}
	qRepo := &mockQuizRepo{
		submissions: []domain.LessonSubmissionDetail{},
	}

	svc := progressService.NewProgressService(pRepo, nil, qRepo, nil)

	prog, err := svc.GetLessonProgress(ctx, 10, 42)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if prog == nil {
		t.Fatal("expected non-nil progress")
	}

	if prog.Status != domain.ProgressStatusNotStarted {
		t.Errorf("expected status %s, got %s", domain.ProgressStatusNotStarted, prog.Status)
	}

	if prog.Score != nil {
		t.Errorf("expected nil score, got %v", *prog.Score)
	}

	if prog.LessonID != 42 {
		t.Errorf("expected lessonID 42, got %d", prog.LessonID)
	}

	if prog.Submissions == nil || len(prog.Submissions) != 0 {
		t.Errorf("expected empty submissions slice, got %v", prog.Submissions)
	}
}

func TestGetLessonProgress_ExistingProgress_WithSubmissions(t *testing.T) {
	ctx := context.Background()
	score := 95
	feedback := "Отличная работа!"

	pRepo := &mockProgressRepo{
		lessonProgress: &domain.LessonProgress{
			ID:        1,
			UserID:    10,
			LessonID:  42,
			Status:    domain.ProgressStatusCompleted,
			Score:     &score,
			TimeSpent: 1200,
			LastPos:   0,
		},
	}

	qRepo := &mockQuizRepo{
		submissions: []domain.LessonSubmissionDetail{
			{
				QuestionText:    "Практическое задание №1",
				StudentAnswer:   "Мое решение задачи",
				PointsAwarded:   25,
				MaxPoints:       25,
				TeacherFeedback: &feedback,
				IsGraded:        true,
			},
		},
	}

	svc := progressService.NewProgressService(pRepo, nil, qRepo, nil)

	prog, err := svc.GetLessonProgress(ctx, 10, 42)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if prog.Status != domain.ProgressStatusCompleted {
		t.Errorf("expected status %s, got %s", domain.ProgressStatusCompleted, prog.Status)
	}

	if prog.Score == nil || *prog.Score != 95 {
		t.Errorf("expected score 95, got %v", prog.Score)
	}

	if len(prog.Submissions) != 1 {
		t.Fatalf("expected 1 submission, got %d", len(prog.Submissions))
	}

	sub := prog.Submissions[0]
	if sub.QuestionText != "Практическое задание №1" {
		t.Errorf("unexpected question text: %s", sub.QuestionText)
	}
	if sub.PointsAwarded != 25 || sub.MaxPoints != 25 {
		t.Errorf("unexpected points: awarded=%d, max=%d", sub.PointsAwarded, sub.MaxPoints)
	}
	if sub.TeacherFeedback == nil || *sub.TeacherFeedback != feedback {
		t.Errorf("unexpected feedback: %v", sub.TeacherFeedback)
	}
	if !sub.IsGraded {
		t.Errorf("expected is_graded = true")
	}
}

func TestGetLessonProgress_PgxErrNoRows_ReturnsDefaultNotStarted(t *testing.T) {
	ctx := context.Background()
	pRepo := &mockProgressRepo{
		getErr: pgx.ErrNoRows,
	}
	qRepo := &mockQuizRepo{
		submissions: []domain.LessonSubmissionDetail{},
	}

	svc := progressService.NewProgressService(pRepo, nil, qRepo, nil)

	prog, err := svc.GetLessonProgress(ctx, 10, 42)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if prog.Status != domain.ProgressStatusNotStarted {
		t.Errorf("expected status %s, got %s", domain.ProgressStatusNotStarted, prog.Status)
	}
}

func TestGetLessonProgress_RepoError_ReturnsError(t *testing.T) {
	ctx := context.Background()
	dbErr := errors.New("connection failed")
	pRepo := &mockProgressRepo{
		getErr: dbErr,
	}

	svc := progressService.NewProgressService(pRepo, nil, nil, nil)

	prog, err := svc.GetLessonProgress(ctx, 10, 42)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if prog != nil {
		t.Fatalf("expected nil progress on error, got %v", prog)
	}
}
