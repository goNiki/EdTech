package quiz

import (
	"context"
	"fmt"
	"time"

	"edtech/internal/domain"
	errorsAPP "edtech/pkg/errors"

	"github.com/jackc/pgx/v5"
)

func (s *service) StartAttempt(ctx context.Context, userID int64, quizID int64) (*domain.QuizAttempt, error) {
	const op = "service.quiz.StartAttempt"

	quiz, err := s.quizRepo.GetQuizByID(ctx, quizID)
	if err != nil {
		return nil, fmt.Errorf("%s: get quiz: %w", op, err)
	}

	var attempt *domain.QuizAttempt

	err = s.txManager.WithTX(ctx, pgx.TxOptions{}, func(ctx context.Context) error {
		// Advisory lock для защиты от phantom reads при проверке лимита попыток
		if lockErr := s.quizRepo.AcquireAdvisoryLock(ctx, userID, quizID); lockErr != nil {
			return lockErr
		}

		count, txErr := s.quizRepo.CountUserAttemptsForUpdate(ctx, userID, quizID)
		if txErr != nil {
			return txErr
		}

		if quiz.MaxAttempts != nil && *quiz.MaxAttempts > 0 && count >= *quiz.MaxAttempts {
			return errorsAPP.ErrMaxAttemptsReached
		}

		newAttempt := &domain.QuizAttempt{
			QuizID:    quizID,
			UserID:    userID,
			Score:     0,
			Passed:    false,
			StartedAt: time.Now(),
		}

		var createErr error
		attempt, createErr = s.quizRepo.CreateAttempt(ctx, newAttempt)
		if createErr != nil {
			return createErr
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return attempt, nil
}
