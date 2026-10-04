package quiz

import (
	"context"
	"fmt"

	"edtech/internal/domain"
	errorsAPP "edtech/pkg/errors"

	"github.com/jackc/pgx/v5"
)

func (s *service) GradeAttemptAnswer(ctx context.Context, teacherID int64, attemptID int64, answerID int64, points int, feedback *string) (*domain.QuizAttempt, error) {
	const op = "service.quiz.GradeAttemptAnswer"

	var attempt *domain.QuizAttempt

	err := s.txManager.WithTX(ctx, pgx.TxOptions{}, func(ctx context.Context) error {
		var txErr error
		attempt, txErr = s.quizRepo.GetAttemptForUpdate(ctx, attemptID)
		if txErr != nil {
			return txErr
		}

		quiz, txErr := s.quizRepo.GetQuizByID(ctx, attempt.QuizID)
		if txErr != nil {
			return txErr
		}

		if txErr = s.checkTeacherAccess(ctx, teacherID, quiz.LessonID); txErr != nil {
			return txErr
		}

		isCorrect := points > 0
		if txErr = s.quizRepo.UpdateAttemptAnswer(ctx, answerID, points, feedback, isCorrect); txErr != nil {
			return txErr
		}

		ungradedCount, txErr := s.quizRepo.CountUngradedAnswers(ctx, attemptID)
		if txErr != nil {
			return txErr
		}

		if ungradedCount == 0 {
			if txErr = s.finalizeGrading(ctx, quiz, attempt); txErr != nil {
				return txErr
			}
		} else {
			if txErr = s.quizRepo.UpdateAttempt(ctx, attempt); txErr != nil {
				return txErr
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return attempt, nil
}

func (s *service) checkTeacherAccess(ctx context.Context, teacherID, lessonID int64) error {
	lesson, err := s.lessonRepo.GetLessonByID(ctx, lessonID)
	if err != nil {
		return fmt.Errorf("get lesson: %w", err)
	}

	course, err := s.courseRepo.GetCourseByID(ctx, lesson.CourseID)
	if err != nil {
		return fmt.Errorf("get course: %w", err)
	}

	canEdit, err := s.accessService.CanEditCourse(ctx, course, teacherID)
	if err != nil {
		return fmt.Errorf("%w: %w", errorsAPP.ErrCheckingPermissions, err)
	}
	if !canEdit {
		return errorsAPP.ErrForbidden
	}

	return nil
}

func (s *service) finalizeGrading(ctx context.Context, quiz *domain.Quiz, attempt *domain.QuizAttempt) error {
	correctPoints, err := s.quizRepo.SumAttemptPoints(ctx, attempt.ID)
	if err != nil {
		return err
	}

	maxPoints, err := s.quizRepo.GetQuizTotalPoints(ctx, attempt.QuizID)
	if err != nil {
		maxPoints = 1
	}

	attempt.CalculateScore(correctPoints, maxPoints, quiz.PassingScor)

	if err = s.quizRepo.UpdateAttempt(ctx, attempt); err != nil {
		return err
	}

	if attempt.Passed {
		_, _ = s.progressService.CompleteLesson(ctx, attempt.UserID, quiz.LessonID, &attempt.Score, nil, nil)
		_ = s.quizRepo.UpdateLessonProgressAfterQuiz(ctx, attempt.UserID, quiz.LessonID, attempt.Score)
	}

	return nil
}
