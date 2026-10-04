package progress_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"edtech/internal/domain"
	"edtech/internal/infrastructure/txmanager"
	"edtech/internal/repository"
	progressService "edtech/internal/service/progress"
	errorsAPP "edtech/pkg/errors"

	"github.com/jackc/pgx/v5"
)

type mockTxManager struct {
	txmanager.TransactionManager
}

func (m *mockTxManager) WithTX(ctx context.Context, _ pgx.TxOptions, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

type fullMockLessonRepo struct {
	repository.LessonRepository
	lesson  *domain.Lesson
	lessons []domain.Lesson
}

func (m *fullMockLessonRepo) GetLessonByID(_ context.Context, _ int64) (*domain.Lesson, error) {
	return m.lesson, nil
}

func (m *fullMockLessonRepo) GetLessonsByCourseID(_ context.Context, _ int64) ([]domain.Lesson, error) {
	return m.lessons, nil
}

type fullMockProgressRepo struct {
	repository.ProgressRepository
	lessonProgress *domain.LessonProgress
	allProgress    []domain.LessonProgress
	savedScore     int
	statusUpdated  domain.ProgressStatus
	upsertedScore  float64
}

func (m *fullMockProgressRepo) GetLessonProgress(_ context.Context, _, _ int64) (*domain.LessonProgress, error) {
	if m.lessonProgress == nil {
		return nil, errorsAPP.ErrLessonProgressNotFound
	}
	return m.lessonProgress, nil
}

func (m *fullMockProgressRepo) UpdateLessonProgressStatus(_ context.Context, _, _ int64, status domain.ProgressStatus) error {
	m.statusUpdated = status
	return nil
}

func (m *fullMockProgressRepo) UpdateLessonProgressScore(_ context.Context, _, _ int64, score int) error {
	m.savedScore = score
	return nil
}

func (m *fullMockProgressRepo) CreateLessonProgress(_ context.Context, lp *domain.LessonProgress) error {
	if lp.Score != nil {
		m.savedScore = *lp.Score
	}
	m.statusUpdated = lp.Status
	return nil
}

func (m *fullMockProgressRepo) GetAllLessonProgressByCourse(_ context.Context, _, _ int64) ([]domain.LessonProgress, error) {
	return m.allProgress, nil
}

func (m *fullMockProgressRepo) UpsertCourseProgressWithScore(_ context.Context, _, _ int64, _, _ int, _ float64, avgScore float64) error {
	m.upsertedScore = avgScore
	return nil
}

type fullMockQuizRepo struct {
	repository.QuizRepository
	quiz        *domain.Quiz
	attempt     *domain.QuizAttempt
	totalPoints int
	earnedPts   int
}

func (m *fullMockQuizRepo) GetQuizByID(_ context.Context, _ int64) (*domain.Quiz, error) {
	if m.quiz == nil {
		return nil, errorsAPP.ErrQuizNotFound
	}
	return m.quiz, nil
}

func (m *fullMockQuizRepo) GetQuizByLessonID(_ context.Context, _ int64) (*domain.Quiz, error) {
	if m.quiz == nil {
		return nil, errorsAPP.ErrQuizNotFound
	}
	return m.quiz, nil
}

func (m *fullMockQuizRepo) GetAttemptByID(_ context.Context, _ int64) (*domain.QuizAttempt, error) {
	if m.attempt == nil {
		return nil, errorsAPP.ErrAttemptNotFound
	}
	return m.attempt, nil
}

func (m *fullMockQuizRepo) UpdateAttempt(_ context.Context, att *domain.QuizAttempt) error {
	m.attempt = att
	return nil
}

func (m *fullMockQuizRepo) SumAttemptPoints(_ context.Context, _ int64) (int, error) {
	return m.earnedPts, nil
}

func (m *fullMockQuizRepo) GetQuizTotalPoints(_ context.Context, _ int64) (int, error) {
	if m.totalPoints == 0 {
		return 1, nil
	}
	return m.totalPoints, nil
}

func (m *fullMockQuizRepo) SaveEssaySubmission(_ context.Context, _, _, _ int64, _ domain.EssaySubmission) error {
	return nil
}

func (m *fullMockQuizRepo) ListUserAttemptsByLessonID(_ context.Context, _, _ int64) ([]domain.QuizAttempt, error) {
	if m.attempt != nil {
		return []domain.QuizAttempt{*m.attempt}, nil
	}
	return nil, nil
}

func (m *fullMockQuizRepo) GetBestScoreByLessonID(_ context.Context, _, _ int64) (int, error) {
	if m.attempt != nil {
		return m.attempt.Score, nil
	}
	return 0, nil
}

func (m *fullMockQuizRepo) CountUserAttempts(_ context.Context, _, _ int64) (int, error) {
	if m.attempt != nil {
		return 1, nil
	}
	return 0, nil
}

func (m *fullMockQuizRepo) CreateAttempt(_ context.Context, att *domain.QuizAttempt) (*domain.QuizAttempt, error) {
	att.ID = 1001
	m.attempt = att
	return att, nil
}

func (m *fullMockQuizRepo) SaveAttemptDraft(_ context.Context, attemptID, _ int64, currentStep int, draftAnswers []byte) error {
	if m.attempt != nil && m.attempt.ID == attemptID {
		m.attempt.CurrentStep = currentStep
		var d map[string]any
		_ = json.Unmarshal(draftAnswers, &d)
		m.attempt.DraftAnswers = d
		return nil
	}
	return nil
}

func (m *fullMockQuizRepo) GetActiveAttempt(_ context.Context, _, _ int64) (*domain.QuizAttempt, error) {
	if m.attempt != nil && m.attempt.CompletedAt == nil {
		return m.attempt, nil
	}
	return nil, nil
}

// Test 1: Обычная лекция без тестов завершается со 100% прогрессом
func TestCompleteLesson_Lecture_Default100(t *testing.T) {
	ctx := context.Background()
	lRepo := &fullMockLessonRepo{
		lesson: &domain.Lesson{
			ID:       1,
			CourseID: 10,
			Type:     "lecture",
			Content:  `{"content":[{"type":"TextBlock","props":{"text":"Привет"}}]}`,
		},
		lessons: []domain.Lesson{{ID: 1, CourseID: 10}},
	}
	pRepo := &fullMockProgressRepo{}
	qRepo := &fullMockQuizRepo{}

	svc := progressService.NewProgressService(pRepo, lRepo, qRepo, &mockTxManager{})

	res, err := svc.CompleteLesson(ctx, 42, 1, domain.CompleteLessonInput{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Score != 100 {
		t.Errorf("expected score 100 for lecture, got %d", res.Score)
	}
	if !res.IsPassed {
		t.Errorf("expected is_passed = true for lecture")
	}
	if pRepo.savedScore != 100 {
		t.Errorf("expected savedScore 100, got %d", pRepo.savedScore)
	}
}

// Test 2: Урок с проверочным тестом (Puck blocks) - клиент передает score: 100 без ответов -> СЕРВЕР ОБНУЛЯЕТ БАЛЛ
func TestCompleteLesson_QuizBlocks_IgnoreClient100Cheat(t *testing.T) {
	ctx := context.Background()
	puckQuizContent := `{
		"content": [
			{
				"type": "QuizSingleBlock",
				"props": {
					"id": "q1",
					"points": 10,
					"options": [
						{"text": "A", "isCorrect": false},
						{"text": "B", "isCorrect": true}
					]
				}
			}
		]
	}`

	lRepo := &fullMockLessonRepo{
		lesson: &domain.Lesson{
			ID:       2,
			CourseID: 10,
			Type:     "lecture", // даже если тип указан lecture, наличие Puck блоков делает его квизом
			Content:  puckQuizContent,
		},
		lessons: []domain.Lesson{{ID: 2, CourseID: 10}},
	}
	pRepo := &fullMockProgressRepo{}
	qRepo := &fullMockQuizRepo{}

	svc := progressService.NewProgressService(pRepo, lRepo, qRepo, &mockTxManager{})

	clientCheatScore := 100
	res, err := svc.CompleteLesson(ctx, 42, 2, domain.CompleteLessonInput{
		Score:   &clientCheatScore,
		Answers: nil, // студент не дал ответов!
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Score != 0 {
		t.Errorf("expected score 0 (client cheat ignored), got %d", res.Score)
	}
	if res.EarnedPoints != 0 {
		t.Errorf("expected earned_points 0, got %d", res.EarnedPoints)
	}
	if res.TotalMaxPoints != 10 {
		t.Errorf("expected total_max_points 10, got %d", res.TotalMaxPoints)
	}
	if res.IsPassed {
		t.Errorf("expected is_passed = false")
	}
	if pRepo.savedScore != 0 {
		t.Errorf("expected pRepo.savedScore 0, got %d", pRepo.savedScore)
	}
}

// Test 3: Урок с Puck тестами - студент отвечает правильно -> получает заслуженные баллы
func TestCompleteLesson_QuizBlocks_CorrectAnswer(t *testing.T) {
	ctx := context.Background()
	puckQuizContent := `{
		"content": [
			{
				"type": "QuizSingleBlock",
				"props": {
					"id": "q1",
					"points": 10,
					"options": [
						{"text": "A", "isCorrect": false},
						{"text": "B", "isCorrect": true}
					]
				}
			}
		]
	}`

	lRepo := &fullMockLessonRepo{
		lesson: &domain.Lesson{
			ID:       3,
			CourseID: 10,
			Type:     "test",
			Content:  puckQuizContent,
		},
		lessons: []domain.Lesson{{ID: 3, CourseID: 10}},
	}
	pRepo := &fullMockProgressRepo{}
	qRepo := &fullMockQuizRepo{}

	svc := progressService.NewProgressService(pRepo, lRepo, qRepo, &mockTxManager{})

	res, err := svc.CompleteLesson(ctx, 42, 3, domain.CompleteLessonInput{
		Answers: []domain.LessonAnswerSubmission{
			{
				BlockID: "q1",
				Answer:  map[string]any{"selected_option": 1}, // ответ "B"
			},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Score != 100 {
		t.Errorf("expected score 100 for correct answer, got %d", res.Score)
	}
	if res.EarnedPoints != 10 {
		t.Errorf("expected earned_points 10, got %d", res.EarnedPoints)
	}
	if !res.IsPassed {
		t.Errorf("expected is_passed = true")
	}
	if pRepo.savedScore != 100 {
		t.Errorf("expected pRepo.savedScore 100, got %d", pRepo.savedScore)
	}
}

// Test 4: Досрочное завершение теста (is_abandoned = true) с несколькими блоками
func TestCompleteLesson_QuizBlocks_Abandoned(t *testing.T) {
	ctx := context.Background()
	puckQuizContent := `{
		"content": [
			{
				"type": "QuizSingleBlock",
				"props": {
					"id": "q1",
					"points": 10,
					"options": [
						{"text": "A", "isCorrect": true},
						{"text": "B", "isCorrect": false}
					]
				}
			},
			{
				"type": "QuizSingleBlock",
				"props": {
					"id": "q2",
					"points": 10,
					"options": [
						{"text": "C", "isCorrect": true},
						{"text": "D", "isCorrect": false}
					]
				}
			}
		]
	}`

	lRepo := &fullMockLessonRepo{
		lesson: &domain.Lesson{
			ID:       4,
			CourseID: 10,
			Type:     "test",
			Content:  puckQuizContent,
		},
		lessons: []domain.Lesson{{ID: 4, CourseID: 10}},
	}
	pRepo := &fullMockProgressRepo{}
	qRepo := &fullMockQuizRepo{}

	svc := progressService.NewProgressService(pRepo, lRepo, qRepo, &mockTxManager{})

	// Студент решил q1, но бросил тест (q2 не отвечен)
	res, err := svc.CompleteLesson(ctx, 42, 4, domain.CompleteLessonInput{
		IsAbandoned: true,
		Answers: []domain.LessonAnswerSubmission{
			{
				BlockID: "q1",
				Answer:  map[string]any{"selected_option": 0},
			},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 10 из 20 баллов = 50%
	if res.Score != 50 {
		t.Errorf("expected score 50, got %d", res.Score)
	}
	if res.EarnedPoints != 10 {
		t.Errorf("expected earned_points 10, got %d", res.EarnedPoints)
	}
	if res.TotalMaxPoints != 20 {
		t.Errorf("expected total_max_points 20, got %d", res.TotalMaxPoints)
	}
	if res.IsPassed {
		t.Errorf("expected is_passed = false (50 < 70)")
	}
	if pRepo.savedScore != 50 {
		t.Errorf("expected pRepo.savedScore 50, got %d", pRepo.savedScore)
	}
}

// Test 5: Relational quiz attempt с проверкой принадлежности пользователю и уроку
func TestCompleteLesson_RelationalAttempt_SecurityAndScoring(t *testing.T) {
	ctx := context.Background()
	lRepo := &fullMockLessonRepo{
		lesson: &domain.Lesson{
			ID:       5,
			CourseID: 10,
			Type:     "test",
		},
		lessons: []domain.Lesson{{ID: 5, CourseID: 10}},
	}

	qRepo := &fullMockQuizRepo{
		quiz: &domain.Quiz{
			ID:          99,
			LessonID:    5,
			PassingScor: 70,
		},
		attempt: &domain.QuizAttempt{
			ID:     123,
			QuizID: 99,
			UserID: 42,
		},
		totalPoints: 20,
		earnedPts:   16, // 16 / 20 = 80%
	}
	pRepo := &fullMockProgressRepo{}

	svc := progressService.NewProgressService(pRepo, lRepo, qRepo, &mockTxManager{})

	attemptID := int64(123)
	res, err := svc.CompleteLesson(ctx, 42, 5, domain.CompleteLessonInput{
		AttemptID:   &attemptID,
		IsAbandoned: false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Score != 80 {
		t.Errorf("expected score 80, got %d", res.Score)
	}
	if !res.IsPassed {
		t.Errorf("expected is_passed = true")
	}
	if qRepo.attempt.CompletedAt == nil {
		t.Errorf("expected attempt.CompletedAt to be set")
	}

	// Попытка чужого пользователя -> ErrForbidden
	foreignAttemptID := int64(123)
	_, fErr := svc.CompleteLesson(ctx, 999, 5, domain.CompleteLessonInput{
		AttemptID: &foreignAttemptID,
	})
	if fErr == nil {
		t.Fatal("expected ErrForbidden for foreign user, got nil")
	}
}

// Test 6: Best Score Preservation - предыдущий высший балл не перезаписывается худшим
func TestCompleteLesson_BestScorePreservation(t *testing.T) {
	ctx := context.Background()
	prevScore := 90
	pRepo := &fullMockProgressRepo{
		lessonProgress: &domain.LessonProgress{
			Score: &prevScore,
		},
	}
	lRepo := &fullMockLessonRepo{
		lesson: &domain.Lesson{
			ID:       6,
			CourseID: 10,
			Type:     "test",
			Content: `{
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
			}`,
		},
		lessons: []domain.Lesson{{ID: 6, CourseID: 10}},
	}
	qRepo := &fullMockQuizRepo{}

	svc := progressService.NewProgressService(pRepo, lRepo, qRepo, &mockTxManager{})

	// Студент провалил попытку (0 баллов)
	res, err := svc.CompleteLesson(ctx, 42, 6, domain.CompleteLessonInput{
		Answers: nil,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Текущий результат этой попытки - 0
	if res.Score != 0 {
		t.Errorf("expected attempt score 0, got %d", res.Score)
	}
	// Но в прогрессе урока сохранен лучший балл: 90!
	if pRepo.savedScore != 90 {
		t.Errorf("expected preserved best score 90, got %d", pRepo.savedScore)
	}
}

// Test 7: Timeout Verification - превышение time_limit_minutes + 15s приводит к timed_out и отклонению попытки
func TestCompleteLesson_TimeoutVerification(t *testing.T) {
	ctx := context.Background()
	startedAt := time.Now().Add(-16 * time.Minute) // Начата 16 минут назад (лимит 15 мин + 15 сек grace period)

	lRepo := &fullMockLessonRepo{
		lesson: &domain.Lesson{
			ID:       7,
			CourseID: 10,
			Type:     "test",
			QuizSettings: &domain.QuizSettings{
				TimeLimitMinutes:    15,
				PassingScorePercent: 70,
			},
			Content: `{"content": [{"type": "QuizSingleBlock", "props": {"id": "q1", "points": 10, "options": [{"text": "A", "isCorrect": true}]}}]}`,
		},
		lessons: []domain.Lesson{{ID: 7, CourseID: 10}},
	}
	qRepo := &fullMockQuizRepo{
		quiz: &domain.Quiz{
			ID:          70,
			LessonID:    7,
			PassingScor: 70,
		},
		attempt: &domain.QuizAttempt{
			ID:        555,
			QuizID:    70,
			UserID:    42,
			StartedAt: startedAt,
		},
		totalPoints: 10,
		earnedPts:   10,
	}
	pRepo := &fullMockProgressRepo{}

	svc := progressService.NewProgressService(pRepo, lRepo, qRepo, &mockTxManager{})

	attID := int64(555)
	res, err := svc.CompleteLesson(ctx, 42, 7, domain.CompleteLessonInput{
		AttemptID: &attID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Status != "timed_out" {
		t.Errorf("expected status 'timed_out', got %q", res.Status)
	}
	if res.IsPassed {
		t.Errorf("expected is_passed = false on timeout")
	}
	if res.Score != 0 {
		t.Errorf("expected score = 0 on timeout, got %d", res.Score)
	}
	if qRepo.attempt.Passed {
		t.Errorf("expected attempt.Passed = false")
	}
	if qRepo.attempt.CompletedAt == nil {
		t.Errorf("expected attempt.CompletedAt to be set")
	}
}

// Test 8: Exam Blind Mode - в режиме exam_blind правильные ответы скрываются
func TestCompleteLesson_ExamBlindMode(t *testing.T) {
	ctx := context.Background()

	lRepo := &fullMockLessonRepo{
		lesson: &domain.Lesson{
			ID:       8,
			CourseID: 10,
			Type:     "test",
			QuizSettings: &domain.QuizSettings{
				FeedbackMode:        "exam_blind",
				PassingScorePercent: 70,
			},
			Content: `{"content": [{"type": "QuizSingleBlock", "props": {"id": "q1", "points": 10, "options": [{"text": "A", "isCorrect": true}]}}]}`,
		},
		lessons: []domain.Lesson{{ID: 8, CourseID: 10}},
	}
	pRepo := &fullMockProgressRepo{}
	qRepo := &fullMockQuizRepo{}

	svc := progressService.NewProgressService(pRepo, lRepo, qRepo, &mockTxManager{})

	res, err := svc.CompleteLesson(ctx, 42, 8, domain.CompleteLessonInput{
		Answers: []domain.LessonAnswerSubmission{
			{
				BlockID: "q1",
				Answer:  map[string]any{"selected_option": 0},
			},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Results != nil {
		t.Errorf("expected Results to be nil in exam_blind mode, got %v", res.Results)
	}
	if res.Score != 100 {
		t.Errorf("expected score 100, got %d", res.Score)
	}
	if !res.IsPassed {
		t.Errorf("expected is_passed = true")
	}
}

// Test 9: Custom Passing Score Percent - порог из quiz_settings учитывается
func TestCompleteLesson_CustomPassingScorePercent(t *testing.T) {
	ctx := context.Background()

	lRepo := &fullMockLessonRepo{
		lesson: &domain.Lesson{
			ID:       9,
			CourseID: 10,
			Type:     "test",
			QuizSettings: &domain.QuizSettings{
				PassingScorePercent: 85, // Порог повышен до 85%
			},
			Content: `{"content": [
				{"type": "QuizSingleBlock", "props": {"id": "q1", "points": 10, "options": [{"text": "A", "isCorrect": true}]}},
				{"type": "QuizSingleBlock", "props": {"id": "q2", "points": 10, "options": [{"text": "B", "isCorrect": true}]}},
				{"type": "QuizSingleBlock", "props": {"id": "q3", "points": 10, "options": [{"text": "C", "isCorrect": true}]}},
				{"type": "QuizSingleBlock", "props": {"id": "q4", "points": 10, "options": [{"text": "D", "isCorrect": true}]}},
				{"type": "QuizSingleBlock", "props": {"id": "q5", "points": 10, "options": [{"text": "E", "isCorrect": true}]}}
			]}`,
		},
		lessons: []domain.Lesson{{ID: 9, CourseID: 10}},
	}
	pRepo := &fullMockProgressRepo{}
	qRepo := &fullMockQuizRepo{}

	svc := progressService.NewProgressService(pRepo, lRepo, qRepo, &mockTxManager{})

	// Студент ответил правильно на 4 из 5 (80%).
	// При дефолтном пороге 70% было бы зачтено, но при 85% должно быть не зачтено!
	res, err := svc.CompleteLesson(ctx, 42, 9, domain.CompleteLessonInput{
		Answers: []domain.LessonAnswerSubmission{
			{BlockID: "q1", Answer: map[string]any{"selected_option": 0}},
			{BlockID: "q2", Answer: map[string]any{"selected_option": 0}},
			{BlockID: "q3", Answer: map[string]any{"selected_option": 0}},
			{BlockID: "q4", Answer: map[string]any{"selected_option": 0}},
			{BlockID: "q5", Answer: map[string]any{"selected_option": 999}}, // неверно
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Score != 80 {
		t.Errorf("expected score 80, got %d", res.Score)
	}
	if res.IsPassed {
		t.Errorf("expected is_passed = false for 80%% score with 85%% passing threshold")
	}
}
