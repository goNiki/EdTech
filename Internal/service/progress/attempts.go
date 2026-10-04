package progress

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"edtech/internal/domain"
	errorsAPP "edtech/pkg/errors"
)

func (s *service) StartLessonAttempt(ctx context.Context, userID, lessonID int64) (*domain.StartAttemptResult, error) {
	const op = "service.progress.StartLessonAttempt"

	lesson, err := s.lessonRepo.GetLessonByID(ctx, lessonID)
	if err != nil {
		return nil, fmt.Errorf("%s: get lesson: %w", op, err)
	}

	quizSettings := lesson.GetQuizSettings()

	if s.quizRepo == nil {
		return nil, fmt.Errorf("%s: quiz repo unavailable", op)
	}

	quiz, qErr := s.quizRepo.GetQuizByLessonID(ctx, lessonID)
	if qErr != nil {
		if errors.Is(qErr, errorsAPP.ErrQuizNotFound) {
			newQuiz := &domain.Quiz{
				LessonID:    lessonID,
				Title:       lesson.Title,
				PassingScor: quizSettings.PassingScorePercent,
			}
			created, cErr := s.quizRepo.CreateQuiz(ctx, newQuiz)
			if cErr != nil {
				return nil, fmt.Errorf("%s: create quiz: %w", op, cErr)
			}
			quiz = created
		} else {
			return nil, fmt.Errorf("%s: get quiz: %w", op, qErr)
		}
	}

	totalAttempts, err := s.quizRepo.CountUserAttempts(ctx, userID, quiz.ID)
	if err != nil {
		return nil, fmt.Errorf("%s: count user attempts: %w", op, err)
	}

	maxAttempts := quizSettings.MaxAttempts
	if maxAttempts <= 0 && quiz.MaxAttempts != nil && *quiz.MaxAttempts > 0 {
		maxAttempts = *quiz.MaxAttempts
	}

	if maxAttempts > 0 && totalAttempts >= maxAttempts {
		return nil, fmt.Errorf("%s: %w", op, errorsAPP.ErrForbidden)
	}

	now := time.Now()
	attempt := &domain.QuizAttempt{
		QuizID:    quiz.ID,
		UserID:    userID,
		Score:     0,
		Passed:    false,
		StartedAt: now,
	}

	created, err := s.quizRepo.CreateAttempt(ctx, attempt)
	if err != nil {
		return nil, fmt.Errorf("%s: create attempt: %w", op, err)
	}

	return &domain.StartAttemptResult{
		AttemptID: created.ID,
		StartedAt: created.StartedAt,
	}, nil
}

func (s *service) GetLessonAttemptsSummary(ctx context.Context, userID, lessonID int64) (*domain.LessonAttemptsSummary, error) {
	const op = "service.progress.GetLessonAttemptsSummary"

	lesson, err := s.lessonRepo.GetLessonByID(ctx, lessonID)
	if err != nil {
		return nil, fmt.Errorf("%s: get lesson: %w", op, err)
	}

	quizSettings := lesson.GetQuizSettings()
	passingThreshold := quizSettings.PassingScorePercent
	maxAttemptsAllowed := quizSettings.MaxAttempts

	if s.quizRepo != nil {
		if quiz, qErr := s.quizRepo.GetQuizByLessonID(ctx, lessonID); qErr == nil && quiz != nil {
			if maxAttemptsAllowed <= 0 && quiz.MaxAttempts != nil {
				maxAttemptsAllowed = *quiz.MaxAttempts
			}
			if lesson.QuizSettings == nil && quiz.PassingScor > 0 {
				passingThreshold = quiz.PassingScor
			}
		}
	}

	var attempts []domain.QuizAttempt
	if s.quizRepo != nil {
		attempts, _ = s.quizRepo.ListUserAttemptsByLessonID(ctx, userID, lessonID)
	}

	var existingProgress *domain.LessonProgress
	if s.progressRepo != nil {
		existingProgress, _ = s.progressRepo.GetLessonProgress(ctx, userID, lessonID)
	}

	totalAttemptsMade := len(attempts)
	bestScore := 0
	var lastAttempt *domain.LessonAttemptItem
	history := make([]domain.LessonAttemptItem, 0, len(attempts))

	for _, att := range attempts {
		if att.CompletedAt != nil {
			item := domain.LessonAttemptItem{
				AttemptID:   att.ID,
				Score:       att.Score,
				SubmittedAt: *att.CompletedAt,
			}
			history = append(history, item)
			if att.Score > bestScore {
				bestScore = att.Score
			}
			itemCopy := item
			lastAttempt = &itemCopy
		}
	}

	if existingProgress != nil && existingProgress.Score != nil && *existingProgress.Score > bestScore {
		bestScore = *existingProgress.Score
	}

	canStartNew := maxAttemptsAllowed <= 0 || totalAttemptsMade < maxAttemptsAllowed

	isPassed := bestScore >= passingThreshold
	if existingProgress != nil && existingProgress.Status == domain.ProgressStatusCompleted {
		isPassed = true
	}

	return &domain.LessonAttemptsSummary{
		LessonID:            lessonID,
		TotalAttemptsMade:   totalAttemptsMade,
		MaxAttemptsAllowed:  maxAttemptsAllowed,
		CanStartNewAttempt:  canStartNew,
		BestScore:           bestScore,
		BestScorePercentage: bestScore,
		IsPassed:            isPassed,
		PassingThreshold:    passingThreshold,
		LastAttempt:         lastAttempt,
		AttemptsHistory:     history,
	}, nil
}

func (s *service) SaveAttemptDraft(ctx context.Context, userID, lessonID, attemptID int64, currentStep int, answers map[string]any) (time.Time, error) {
	const op = "service.progress.SaveAttemptDraft"

	lesson, err := s.lessonRepo.GetLessonByID(ctx, lessonID)
	if err != nil {
		return time.Time{}, fmt.Errorf("%s: get lesson: %w", op, err)
	}

	if s.quizRepo == nil {
		return time.Time{}, fmt.Errorf("%s: quiz repo unavailable", op)
	}

	attempt, err := s.quizRepo.GetAttemptByID(ctx, attemptID)
	if err != nil {
		return time.Time{}, fmt.Errorf("%s: get attempt: %w", op, err)
	}
	if attempt.UserID != userID {
		return time.Time{}, fmt.Errorf("%s: %w", op, errorsAPP.ErrForbidden)
	}
	if attempt.CompletedAt != nil {
		return time.Time{}, fmt.Errorf("%s: attempt already completed: %w", op, errorsAPP.ErrForbidden)
	}

	quizSettings := lesson.GetQuizSettings()
	if quizSettings.TimeLimitMinutes > 0 {
		gracePeriod := 15 * time.Second
		timeLimit := time.Duration(quizSettings.TimeLimitMinutes)*time.Minute + gracePeriod
		if time.Since(attempt.StartedAt) > timeLimit {
			now := time.Now()
			attempt.CompletedAt = &now
			attempt.Passed = false
			attempt.Score = 0
			_ = s.quizRepo.UpdateAttempt(ctx, attempt)
			return time.Time{}, fmt.Errorf("%s: attempt time expired: %w", op, errorsAPP.ErrForbidden)
		}
	}

	if currentStep <= 0 {
		currentStep = 1
	}

	draftBytes, err := json.Marshal(answers)
	if err != nil {
		return time.Time{}, fmt.Errorf("%s: marshal draft answers: %w", op, err)
	}

	if err := s.quizRepo.SaveAttemptDraft(ctx, attemptID, userID, currentStep, draftBytes); err != nil {
		return time.Time{}, fmt.Errorf("%s: save draft: %w", op, err)
	}

	return time.Now(), nil
}

func (s *service) GetActiveLessonAttempt(ctx context.Context, userID, lessonID int64) (*domain.ActiveAttemptResult, error) {
	const op = "service.progress.GetActiveLessonAttempt"

	lesson, err := s.lessonRepo.GetLessonByID(ctx, lessonID)
	if err != nil {
		return nil, fmt.Errorf("%s: get lesson: %w", op, err)
	}

	if s.quizRepo == nil {
		return nil, fmt.Errorf("%s: quiz repo unavailable", op)
	}

	attempt, err := s.quizRepo.GetActiveAttempt(ctx, userID, lessonID)
	if err != nil {
		return nil, fmt.Errorf("%s: get active attempt: %w", op, err)
	}
	if attempt == nil {
		return &domain.ActiveAttemptResult{
			HasActiveAttempt: false,
		}, nil
	}

	quizSettings := lesson.GetQuizSettings()
	timeLimitMinutes := quizSettings.TimeLimitMinutes
	remainingSeconds := 0

	if timeLimitMinutes > 0 {
		totalLimitSeconds := timeLimitMinutes * 60
		elapsedSeconds := int(time.Since(attempt.StartedAt).Seconds())
		gracePeriodSeconds := 15

		if elapsedSeconds > (totalLimitSeconds + gracePeriodSeconds) {
			// Попытка с истекшим временем автоматически помечается как timed_out и не возвращается как активная
			now := time.Now()
			attempt.CompletedAt = &now
			attempt.Passed = false
			attempt.Score = 0
			_ = s.quizRepo.UpdateAttempt(ctx, attempt)

			return &domain.ActiveAttemptResult{
				HasActiveAttempt: false,
			}, nil
		}

		remainingSeconds = totalLimitSeconds - elapsedSeconds
		if remainingSeconds < 0 {
			remainingSeconds = 0
		}
	}

	return &domain.ActiveAttemptResult{
		HasActiveAttempt: true,
		Attempt: &domain.ActiveAttemptInfo{
			ID:               attempt.ID,
			StartedAt:        attempt.StartedAt,
			TimeLimitMinutes: timeLimitMinutes,
			RemainingSeconds: remainingSeconds,
			CurrentStep:      attempt.CurrentStep,
			DraftAnswers:     attempt.DraftAnswers,
		},
	}, nil
}
