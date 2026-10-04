package quiz

import (
	"context"
	"fmt"
	"time"

	"edtech/internal/domain"
	errorsAPP "edtech/pkg/errors"

	"github.com/jackc/pgx/v5"
)

func (s *service) SubmitAttempt(ctx context.Context, userID int64, attemptID int64, answers []domain.QuizAttemptAnswer) (*domain.QuizAttempt, error) {
	const op = "service.quiz.SubmitAttempt"

	var attempt *domain.QuizAttempt

	err := s.txManager.WithTX(ctx, pgx.TxOptions{}, func(ctx context.Context) error {
		var txErr error
		attempt, txErr = s.quizRepo.GetAttemptForUpdate(ctx, attemptID)
		if txErr != nil {
			return txErr
		}

		if attempt.UserID != userID {
			return errorsAPP.ErrForbidden
		}

		if attempt.CompletedAt != nil {
			return errorsAPP.ErrAttemptAlreadyCompleted
		}

		quiz, txErr := s.quizRepo.GetQuizByID(ctx, attempt.QuizID)
		if txErr != nil {
			return txErr
		}

		if txErr = s.checkTimeLimit(quiz, attempt); txErr != nil {
			return txErr
		}

		hasOpenText := s.autoGradeAnswers(ctx, attemptID, answers)

		if len(answers) > 0 {
			if txErr = s.quizRepo.CreateBatchAnswers(ctx, answers); txErr != nil {
				return txErr
			}
		}

		now := time.Now()
		attempt.CompletedAt = &now

		if hasOpenText {
			attempt.NeedsGrading = true
			attempt.Passed = false
			attempt.Score = 0
		} else {
			if txErr = s.finalizeAttemptScore(ctx, quiz, attempt, answers); txErr != nil {
				return txErr
			}
		}

		return s.quizRepo.UpdateAttempt(ctx, attempt)
	})

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return attempt, nil
}

func (s *service) checkTimeLimit(quiz *domain.Quiz, attempt *domain.QuizAttempt) error {
	if quiz.TimeLimit != nil && *quiz.TimeLimit > 0 {
		timeLimitDuration := time.Duration(*quiz.TimeLimit+30) * time.Second
		if time.Since(attempt.StartedAt) > timeLimitDuration {
			return errorsAPP.ErrTimeLimitExceeded
		}
	}
	return nil
}

func (s *service) autoGradeAnswers(ctx context.Context, attemptID int64, answers []domain.QuizAttemptAnswer) bool {
	hasOpenText := false
	for i := range answers {
		answers[i].AttemptID = attemptID

		if answers[i].AnswerID == nil {
			hasOpenText = true
			answers[i].IsCorrect = nil
			answers[i].Points = 0
		} else {
			isCorrect, points, err := s.quizRepo.GetAnswerPointsAndCorrectness(ctx, *answers[i].AnswerID)
			if err != nil {
				isCorrect = false
				points = 0
			}
			answers[i].IsCorrect = &isCorrect
			if isCorrect {
				answers[i].Points = points
			} else {
				answers[i].Points = 0
			}
		}
	}
	return hasOpenText
}

func (s *service) finalizeAttemptScore(ctx context.Context, quiz *domain.Quiz, attempt *domain.QuizAttempt, answers []domain.QuizAttemptAnswer) error {
	correctPoints := 0
	for _, ans := range answers {
		correctPoints += ans.Points
	}

	maxPoints, err := s.quizRepo.GetQuizTotalPoints(ctx, attempt.QuizID)
	if err != nil {
		maxPoints = 1
	}

	attempt.CalculateScore(correctPoints, maxPoints, quiz.PassingScor)

	if attempt.Passed {
		if _, err := s.progressService.CompleteLesson(ctx, attempt.UserID, quiz.LessonID, domain.CompleteLessonInput{
			Score:     &attempt.Score,
			AttemptID: &attempt.ID,
		}); err != nil {
			return fmt.Errorf("complete lesson: %w", err)
		}
		if err := s.quizRepo.UpdateLessonProgressAfterQuiz(ctx, attempt.UserID, quiz.LessonID, attempt.Score); err != nil {
			return fmt.Errorf("update lesson progress after quiz: %w", err)
		}
	}

	return nil
}
