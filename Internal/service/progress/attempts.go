package progress

import (
	"context"
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

	if s.quizRepo == nil {
		return nil, fmt.Errorf("%s: quiz repo unavailable", op)
	}

	quiz, qErr := s.quizRepo.GetQuizByLessonID(ctx, lessonID)
	if qErr != nil {
		if errors.Is(qErr, errorsAPP.ErrQuizNotFound) {
			newQuiz := &domain.Quiz{
				LessonID:    lessonID,
				Title:       lesson.Title,
				PassingScor: 70,
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

	if quiz.MaxAttempts != nil && *quiz.MaxAttempts > 0 && totalAttempts >= *quiz.MaxAttempts {
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

	_, err := s.lessonRepo.GetLessonByID(ctx, lessonID)
	if err != nil {
		return nil, fmt.Errorf("%s: get lesson: %w", op, err)
	}

	passingThreshold := 70
	maxAttemptsAllowed := 0

	if s.quizRepo != nil {
		if quiz, qErr := s.quizRepo.GetQuizByLessonID(ctx, lessonID); qErr == nil && quiz != nil {
			if quiz.PassingScor > 0 {
				passingThreshold = quiz.PassingScor
			}
			if quiz.MaxAttempts != nil {
				maxAttemptsAllowed = *quiz.MaxAttempts
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
