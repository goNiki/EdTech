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

// 7. Проверка лимита попыток max_attempts из QuizSettings урока
func TestStartAttempt_MaxAttemptsFromQuizSettings(t *testing.T) {
	ctx := context.Background()
	lRepo := &fullMockLessonRepo{
		lesson: &domain.Lesson{
			ID:       20,
			Title:    "Контрольная работа",
			QuizSettings: &domain.QuizSettings{
				MaxAttempts:          2,
				PassingScorePercent:  80,
				TimeLimitMinutes:      20,
			},
		},
	}
	qRepo := &attemptsMockQuizRepo{
		fullMockQuizRepo: fullMockQuizRepo{
			quiz: &domain.Quiz{
				ID:       15,
				LessonID: 20,
			},
		},
		attemptsCount: 2, // Лимит 2 исчерпан
	}
	pRepo := &fullMockProgressRepo{}

	svc := progressService.NewProgressService(pRepo, lRepo, qRepo, &mockTxManager{})

	// Попытка запуска третьей попытки должна быть отклонена
	_, err := svc.StartLessonAttempt(ctx, 42, 20)
	if err == nil {
		t.Fatalf("expected ErrForbidden for exceeding max_attempts, got nil")
	}
	if !errors.Is(err, errorsAPP.ErrForbidden) {
		t.Errorf("expected ErrForbidden, got %v", err)
	}

	// А если сделана 1 попытка из 2 - разрешено
	qRepo.attemptsCount = 1
	startRes, err := svc.StartLessonAttempt(ctx, 42, 20)
	if err != nil {
		t.Fatalf("unexpected error when under limit: %v", err)
	}
	if startRes.AttemptID != 777 {
		t.Errorf("expected attemptID 777, got %d", startRes.AttemptID)
	}
}

// 8. Проверка summary правил тестирования из QuizSettings урока
func TestGetLessonAttemptsSummary_RulesFromQuizSettings(t *testing.T) {
	ctx := context.Background()
	lRepo := &fullMockLessonRepo{
		lesson: &domain.Lesson{
			ID: 30,
			QuizSettings: &domain.QuizSettings{
				MaxAttempts:          1,  // Экзаменационный срез
				PassingScorePercent:  85, // Повышенный порог
				FeedbackMode:         "exam_blind",
				TimeLimitMinutes:     30,
			},
		},
	}
	qRepo := &attemptsMockQuizRepo{
		attemptsCount: 0,
		attemptsList:  []domain.QuizAttempt{},
	}
	pRepo := &fullMockProgressRepo{}

	svc := progressService.NewProgressService(pRepo, lRepo, qRepo, &mockTxManager{})

	summary, err := svc.GetLessonAttemptsSummary(ctx, 42, 30)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if summary.MaxAttemptsAllowed != 1 {
		t.Errorf("expected maxAttemptsAllowed 1, got %d", summary.MaxAttemptsAllowed)
	}
	if summary.PassingThreshold != 85 {
		t.Errorf("expected passingThreshold 85, got %d", summary.PassingThreshold)
	}
	if !summary.CanStartNewAttempt {
		t.Errorf("expected canStartNewAttempt true")
	}
}

// 9. Автосохранение драфта ответов и текущего шага
func TestSaveAttemptDraft_Success(t *testing.T) {
	ctx := context.Background()
	lRepo := &fullMockLessonRepo{
		lesson: &domain.Lesson{
			ID: 40,
			QuizSettings: &domain.QuizSettings{
				TimeLimitMinutes: 20,
			},
		},
	}
	qRepo := &attemptsMockQuizRepo{
		fullMockQuizRepo: fullMockQuizRepo{
			attempt: &domain.QuizAttempt{
				ID:        100,
				UserID:    42,
				StartedAt: time.Now(),
			},
		},
	}
	pRepo := &fullMockProgressRepo{}

	svc := progressService.NewProgressService(pRepo, lRepo, qRepo, &mockTxManager{})

	answers := map[string]any{
		"QuizSingleBlock-1": float64(2),
		"QuizMultiBlock-2":  []any{float64(0), float64(3)},
	}

	savedAt, err := svc.SaveAttemptDraft(ctx, 42, 40, 100, 3, answers)
	if err != nil {
		t.Fatalf("unexpected error saving draft: %v", err)
	}
	if savedAt.IsZero() {
		t.Errorf("expected non-zero savedAt timestamp")
	}
	if qRepo.attempt.CurrentStep != 3 {
		t.Errorf("expected current step 3, got %d", qRepo.attempt.CurrentStep)
	}
	if len(qRepo.attempt.DraftAnswers) != 2 {
		t.Errorf("expected 2 answers in draft, got %d", len(qRepo.attempt.DraftAnswers))
	}
}

// 10. Отклонение сохранения драфта для чужого пользователя или завершенной попытки
func TestSaveAttemptDraft_SecurityGuards(t *testing.T) {
	ctx := context.Background()
	completedTime := time.Now()
	lRepo := &fullMockLessonRepo{
		lesson: &domain.Lesson{ID: 40},
	}
	qRepo := &attemptsMockQuizRepo{
		fullMockQuizRepo: fullMockQuizRepo{
			attempt: &domain.QuizAttempt{
				ID:        100,
				UserID:    42,
				StartedAt: time.Now(),
			},
		},
	}
	pRepo := &fullMockProgressRepo{}

	svc := progressService.NewProgressService(pRepo, lRepo, qRepo, &mockTxManager{})

	// 1. Чужой пользователь
	_, err := svc.SaveAttemptDraft(ctx, 999, 40, 100, 2, map[string]any{"q1": 1})
	if err == nil {
		t.Fatalf("expected ErrForbidden for foreign user, got nil")
	}

	// 2. Уже завершенная попытка
	qRepo.attempt.CompletedAt = &completedTime
	_, err = svc.SaveAttemptDraft(ctx, 42, 40, 100, 2, map[string]any{"q1": 1})
	if err == nil {
		t.Fatalf("expected ErrForbidden for completed attempt, got nil")
	}
}

// 11. Отклонение сохранения драфта при истечении времени теста
func TestSaveAttemptDraft_ExpiredTimeout(t *testing.T) {
	ctx := context.Background()
	startedAt := time.Now().Add(-25 * time.Minute) // Лимит 20 минут, прошло 25

	lRepo := &fullMockLessonRepo{
		lesson: &domain.Lesson{
			ID: 40,
			QuizSettings: &domain.QuizSettings{
				TimeLimitMinutes: 20,
			},
		},
	}
	qRepo := &attemptsMockQuizRepo{
		fullMockQuizRepo: fullMockQuizRepo{
			attempt: &domain.QuizAttempt{
				ID:        100,
				UserID:    42,
				StartedAt: startedAt,
			},
		},
	}
	pRepo := &fullMockProgressRepo{}

	svc := progressService.NewProgressService(pRepo, lRepo, qRepo, &mockTxManager{})

	_, err := svc.SaveAttemptDraft(ctx, 42, 40, 100, 2, map[string]any{"q1": 1})
	if err == nil {
		t.Fatalf("expected error for expired attempt, got nil")
	}
	if qRepo.attempt.CompletedAt == nil {
		t.Errorf("expected attempt to be closed with completed_at")
	}
	if qRepo.attempt.Passed {
		t.Errorf("expected passed = false for timed out attempt")
	}
}

// 12. Проверка активной попытки и расчет remaining_seconds
func TestGetActiveLessonAttempt_ActiveWithRemainingSeconds(t *testing.T) {
	ctx := context.Background()
	startedAt := time.Now().Add(-5 * time.Minute) // 5 минут назад, лимит 20 минут -> осталось ~15 минут (900 сек)

	lRepo := &fullMockLessonRepo{
		lesson: &domain.Lesson{
			ID: 50,
			QuizSettings: &domain.QuizSettings{
				TimeLimitMinutes: 20,
			},
		},
	}
	qRepo := &attemptsMockQuizRepo{
		fullMockQuizRepo: fullMockQuizRepo{
			attempt: &domain.QuizAttempt{
				ID:          200,
				UserID:      42,
				StartedAt:   startedAt,
				CurrentStep: 2,
				DraftAnswers: map[string]any{
					"q1": float64(1),
				},
			},
		},
	}
	pRepo := &fullMockProgressRepo{}

	svc := progressService.NewProgressService(pRepo, lRepo, qRepo, &mockTxManager{})

	res, err := svc.GetActiveLessonAttempt(ctx, 42, 50)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !res.HasActiveAttempt {
		t.Fatalf("expected has_active_attempt = true")
	}
	if res.Attempt.ID != 200 {
		t.Errorf("expected attempt ID 200, got %d", res.Attempt.ID)
	}
	if res.Attempt.CurrentStep != 2 {
		t.Errorf("expected current step 2, got %d", res.Attempt.CurrentStep)
	}
	if res.Attempt.RemainingSeconds < 880 || res.Attempt.RemainingSeconds > 910 {
		t.Errorf("expected remaining_seconds around 900, got %d", res.Attempt.RemainingSeconds)
	}
	if len(res.Attempt.DraftAnswers) != 1 {
		t.Errorf("expected 1 draft answer, got %d", len(res.Attempt.DraftAnswers))
	}
}

// 13. Активная попытка с истекшим временем автозакрывается и возвращает false
func TestGetActiveLessonAttempt_ExpiredAutoClosed(t *testing.T) {
	ctx := context.Background()
	startedAt := time.Now().Add(-25 * time.Minute) // 25 минут назад при лимите 20 минут

	lRepo := &fullMockLessonRepo{
		lesson: &domain.Lesson{
			ID: 50,
			QuizSettings: &domain.QuizSettings{
				TimeLimitMinutes: 20,
			},
		},
	}
	qRepo := &attemptsMockQuizRepo{
		fullMockQuizRepo: fullMockQuizRepo{
			attempt: &domain.QuizAttempt{
				ID:        200,
				UserID:    42,
				StartedAt: startedAt,
			},
		},
	}
	pRepo := &fullMockProgressRepo{}

	svc := progressService.NewProgressService(pRepo, lRepo, qRepo, &mockTxManager{})

	res, err := svc.GetActiveLessonAttempt(ctx, 42, 50)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.HasActiveAttempt {
		t.Fatalf("expected has_active_attempt = false for expired attempt")
	}
	if qRepo.attempt.CompletedAt == nil {
		t.Errorf("expected attempt to be closed")
	}
	if qRepo.attempt.Passed {
		t.Errorf("expected passed = false")
	}
}
