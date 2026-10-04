package progress_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"edtech/internal/domain"
	progressService "edtech/internal/service/progress"
	errorsAPP "edtech/pkg/errors"
)

type attemptsMockQuizRepo struct {
	fullMockQuizRepo
	attemptsCount int
	attemptsList  []domain.QuizAttempt
	bestScore     int
	createdAtt    *domain.QuizAttempt
}

func (m *attemptsMockQuizRepo) CountUserAttempts(_ context.Context, _, _ int64) (int, error) {
	return m.attemptsCount, nil
}

func (m *attemptsMockQuizRepo) ListUserAttemptsByLessonID(_ context.Context, _, _ int64) ([]domain.QuizAttempt, error) {
	return m.attemptsList, nil
}

func (m *attemptsMockQuizRepo) GetBestScoreByLessonID(_ context.Context, _, _ int64) (int, error) {
	return m.bestScore, nil
}

func (m *attemptsMockQuizRepo) CreateAttempt(_ context.Context, att *domain.QuizAttempt) (*domain.QuizAttempt, error) {
	att.ID = 777
	m.createdAtt = att
	return att, nil
}

// 1. Pre-flight summary: пустая история попыток
func TestGetLessonAttemptsSummary_Empty(t *testing.T) {
	ctx := context.Background()
	maxAttempts := 3
	lRepo := &fullMockLessonRepo{
		lesson: &domain.Lesson{ID: 10, Title: "Тест по Go"},
	}
	qRepo := &attemptsMockQuizRepo{
		fullMockQuizRepo: fullMockQuizRepo{
			quiz: &domain.Quiz{
				ID:          5,
				LessonID:    10,
				PassingScor: 75,
				MaxAttempts: &maxAttempts,
			},
		},
		attemptsList:  []domain.QuizAttempt{},
		attemptsCount: 0,
	}
	pRepo := &fullMockProgressRepo{}

	svc := progressService.NewProgressService(pRepo, lRepo, qRepo, &mockTxManager{})

	summary, err := svc.GetLessonAttemptsSummary(ctx, 42, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if summary.LessonID != 10 {
		t.Errorf("expected lessonID 10, got %d", summary.LessonID)
	}
	if summary.TotalAttemptsMade != 0 {
		t.Errorf("expected 0 attempts made, got %d", summary.TotalAttemptsMade)
	}
	if summary.MaxAttemptsAllowed != 3 {
		t.Errorf("expected maxAttemptsAllowed 3, got %d", summary.MaxAttemptsAllowed)
	}
	if !summary.CanStartNewAttempt {
		t.Errorf("expected canStartNewAttempt = true")
	}
	if summary.BestScore != 0 {
		t.Errorf("expected bestScore 0, got %d", summary.BestScore)
	}
	if summary.IsPassed {
		t.Errorf("expected isPassed = false")
	}
	if summary.PassingThreshold != 75 {
		t.Errorf("expected passingThreshold 75, got %d", summary.PassingThreshold)
	}
	if summary.LastAttempt != nil {
		t.Errorf("expected nil lastAttempt, got %v", summary.LastAttempt)
	}
	if len(summary.AttemptsHistory) != 0 {
		t.Errorf("expected empty attemptsHistory, got %d", len(summary.AttemptsHistory))
	}
}

// 2. Pre-flight summary: история с несколькими попытками и расчетом лучшего балла
func TestGetLessonAttemptsSummary_WithHistoryAndBestScore(t *testing.T) {
	ctx := context.Background()
	maxAttempts := 3
	t1 := time.Date(2026, 10, 3, 15, 20, 0, 0, time.UTC)
	t2 := time.Date(2026, 10, 4, 10, 30, 0, 0, time.UTC)

	lRepo := &fullMockLessonRepo{
		lesson: &domain.Lesson{ID: 10, Title: "Тест по Go"},
	}
	qRepo := &attemptsMockQuizRepo{
		fullMockQuizRepo: fullMockQuizRepo{
			quiz: &domain.Quiz{
				ID:          5,
				LessonID:    10,
				PassingScor: 70,
				MaxAttempts: &maxAttempts,
			},
		},
		attemptsList: []domain.QuizAttempt{
			{
				ID:          12,
				QuizID:      5,
				UserID:      42,
				Score:       90,
				Passed:      true,
				StartedAt:   t1.Add(-10 * time.Minute),
				CompletedAt: &t1,
			},
			{
				ID:          15,
				QuizID:      5,
				UserID:      42,
				Score:       85,
				Passed:      true,
				StartedAt:   t2.Add(-10 * time.Minute),
				CompletedAt: &t2,
			},
		},
		attemptsCount: 2,
		bestScore:     90,
	}
	pRepo := &fullMockProgressRepo{}

	svc := progressService.NewProgressService(pRepo, lRepo, qRepo, &mockTxManager{})

	summary, err := svc.GetLessonAttemptsSummary(ctx, 42, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if summary.TotalAttemptsMade != 2 {
		t.Errorf("expected totalAttemptsMade 2, got %d", summary.TotalAttemptsMade)
	}
	if summary.BestScore != 90 {
		t.Errorf("expected bestScore 90, got %d", summary.BestScore)
	}
	if summary.BestScorePercentage != 90 {
		t.Errorf("expected bestScorePercentage 90, got %d", summary.BestScorePercentage)
	}
	if !summary.IsPassed {
		t.Errorf("expected isPassed = true")
	}
	if !summary.CanStartNewAttempt {
		t.Errorf("expected canStartNewAttempt = true (2 < 3)")
	}
	if summary.LastAttempt == nil || summary.LastAttempt.AttemptID != 15 || summary.LastAttempt.Score != 85 {
		t.Errorf("unexpected lastAttempt: %+v", summary.LastAttempt)
	}
	if len(summary.AttemptsHistory) != 2 {
		t.Fatalf("expected 2 history items, got %d", len(summary.AttemptsHistory))
	}
	if summary.AttemptsHistory[0].AttemptID != 12 || summary.AttemptsHistory[0].Score != 90 {
		t.Errorf("unexpected first history item: %+v", summary.AttemptsHistory[0])
	}
}

// 3. Старт попытки: успешное создание
func TestStartLessonAttempt_Success(t *testing.T) {
	ctx := context.Background()
	maxAttempts := 3
	lRepo := &fullMockLessonRepo{
		lesson: &domain.Lesson{ID: 10, Title: "Тест по Go"},
	}
	qRepo := &attemptsMockQuizRepo{
		fullMockQuizRepo: fullMockQuizRepo{
			quiz: &domain.Quiz{
				ID:          5,
				LessonID:    10,
				PassingScor: 70,
				MaxAttempts: &maxAttempts,
			},
		},
		attemptsCount: 1, // 1 < 3
	}
	pRepo := &fullMockProgressRepo{}

	svc := progressService.NewProgressService(pRepo, lRepo, qRepo, &mockTxManager{})

	res, err := svc.StartLessonAttempt(ctx, 42, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.AttemptID != 777 {
		t.Errorf("expected attemptID 777, got %d", res.AttemptID)
	}
	if res.StartedAt.IsZero() {
		t.Errorf("expected non-zero startedAt")
	}
}

// 4. Старт попытки: блокировка при исчерпании лимита (403 Forbidden)
func TestStartLessonAttempt_LimitReached(t *testing.T) {
	ctx := context.Background()
	maxAttempts := 2
	lRepo := &fullMockLessonRepo{
		lesson: &domain.Lesson{ID: 10, Title: "Тест по Go"},
	}
	qRepo := &attemptsMockQuizRepo{
		fullMockQuizRepo: fullMockQuizRepo{
			quiz: &domain.Quiz{
				ID:          5,
				LessonID:    10,
				PassingScor: 70,
				MaxAttempts: &maxAttempts,
			},
		},
		attemptsCount: 2, // 2 >= 2 лимит исчерпан
	}
	pRepo := &fullMockProgressRepo{}

	svc := progressService.NewProgressService(pRepo, lRepo, qRepo, &mockTxManager{})

	_, err := svc.StartLessonAttempt(ctx, 42, 10)
	if err == nil {
		t.Fatal("expected 403 Forbidden error, got nil")
	}
	if !errors.Is(err, errorsAPP.ErrForbidden) {
		t.Errorf("expected ErrForbidden, got %v", err)
	}
}

// 5. Best Score Preservation: при прохождении теста на более низкий балл итоговая оценка в уроке не уменьшается
func TestCompleteLesson_BestScorePreservedOverWorseAttempt(t *testing.T) {
	ctx := context.Background()
	lRepo := &fullMockLessonRepo{
		lesson: &domain.Lesson{
			ID:       10,
			CourseID: 100,
			Type:     "test",
		},
		lessons: []domain.Lesson{{ID: 10, CourseID: 100}},
	}
	qRepo := &attemptsMockQuizRepo{
		fullMockQuizRepo: fullMockQuizRepo{
			quiz: &domain.Quiz{
				ID:          5,
				LessonID:    10,
				PassingScor: 70,
			},
			attempt: &domain.QuizAttempt{
				ID:     20,
				QuizID: 5,
				UserID: 42,
			},
			totalPoints: 10,
			earnedPts:   4, // 4 / 10 = 40% (провал попытки)
		},
		bestScore: 90, // в базе уже зафиксирована предыдущая успешная попытка на 90%
	}
	prevScore := 90
	pRepo := &fullMockProgressRepo{
		lessonProgress: &domain.LessonProgress{
			Score:  &prevScore,
			Status: domain.ProgressStatusCompleted,
		},
	}

	svc := progressService.NewProgressService(pRepo, lRepo, qRepo, &mockTxManager{})

	attemptID := int64(20)
	res, err := svc.CompleteLesson(ctx, 42, 10, domain.CompleteLessonInput{
		AttemptID: &attemptID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Результат текущей попытки — 40
	if res.Score != 40 {
		t.Errorf("expected attempt score 40, got %d", res.Score)
	}
	// Но в прогрессе урока сохранена лучшая оценка: 90!
	if pRepo.savedScore != 90 {
		t.Errorf("expected preserved score 90, got %d", pRepo.savedScore)
	}
	// И статус урока НЕ сбросился, остался completed!
	if pRepo.statusUpdated != domain.ProgressStatusCompleted {
		t.Errorf("expected status completed preserved, got %s", pRepo.statusUpdated)
	}
}

// 6. Заброшенная попытка не сбрасывает статус ранее успешно сданного урока
func TestCompleteLesson_AbandonedAttemptDoesNotClearCompletedStatus(t *testing.T) {
	ctx := context.Background()
	puckQuizContent := `{
		"content": [
			{
				"type": "QuizSingleBlock",
				"props": {
					"id": "q1",
					"points": 10,
					"options": [{"text": "A", "isCorrect": true}]
				}
			}
		]
	}`

	lRepo := &fullMockLessonRepo{
		lesson: &domain.Lesson{
			ID:       10,
			CourseID: 100,
			Type:     "test",
			Content:  puckQuizContent,
		},
		lessons: []domain.Lesson{{ID: 10, CourseID: 100}},
	}
	prevScore := 100
	pRepo := &fullMockProgressRepo{
		lessonProgress: &domain.LessonProgress{
			Score:  &prevScore,
			Status: domain.ProgressStatusCompleted, // ранее уже был успешно сдан!
		},
	}
	qRepo := &attemptsMockQuizRepo{}

	svc := progressService.NewProgressService(pRepo, lRepo, qRepo, &mockTxManager{})

	// Студент забросил тест (0 ответов, is_abandoned = true)
	res, err := svc.CompleteLesson(ctx, 42, 10, domain.CompleteLessonInput{
		IsAbandoned: true,
		Answers:     nil,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Score != 0 {
		t.Errorf("expected attempt score 0, got %d", res.Score)
	}
	// Оценка урока осталась 100
	if pRepo.savedScore != 100 {
		t.Errorf("expected savedScore 100, got %d", pRepo.savedScore)
	}
	// Статус остался completed
	if pRepo.statusUpdated != domain.ProgressStatusCompleted {
		t.Errorf("expected status completed preserved, got %s", pRepo.statusUpdated)
	}
}
