package progress

import (
	"context"
	"errors"
	"fmt"
	"time"

	"edtech/internal/domain"
	"edtech/internal/service/quiz"
	errorsAPP "edtech/pkg/errors"

	"github.com/jackc/pgx/v5"
)

type scoreResolution struct {
	finalScore   int
	earnedPoints int
	totalPoints  int
	isPassed     bool
}

func (s *service) CompleteLesson(ctx context.Context, userID int64, lessonID int64, input domain.CompleteLessonInput) (*domain.LessonCompletionResult, error) {
	const op = "service.progress.CompleteLesson"

	lesson, err := s.lessonRepo.GetLessonByID(ctx, lessonID)
	if err != nil {
		return nil, fmt.Errorf("%s: get lesson: %w", op, err)
	}

	validationResult, vErr := quiz.ValidateQuizSubmission(lesson.Content, input.Answers)
	if vErr != nil {
		return nil, fmt.Errorf("%s: validate quiz answers: %w", op, vErr)
	}

	scoreRes, rErr := s.resolveLessonScore(ctx, userID, lessonID, lesson, input, validationResult)
	if rErr != nil {
		return nil, fmt.Errorf("%s: %w", op, rErr)
	}

	// Best Score Preservation: check completed quiz attempts and previous progress score
	progressScore := scoreRes.finalScore
	if s.quizRepo != nil {
		if bestAttemptScore, bErr := s.quizRepo.GetBestScoreByLessonID(ctx, userID, lessonID); bErr == nil && bestAttemptScore > progressScore {
			progressScore = bestAttemptScore
		}
	}

	var existingProgress *domain.LessonProgress
	if s.progressRepo != nil {
		if ep, pErr := s.progressRepo.GetLessonProgress(ctx, userID, lessonID); pErr == nil && ep != nil {
			existingProgress = ep
			if ep.Score != nil && *ep.Score > progressScore {
				progressScore = *ep.Score
			}
			if ep.Status == domain.ProgressStatusCompleted {
				scoreRes.isPassed = true
			}
		}
	}

	targetStatus := domain.ProgressStatusCompleted
	if !scoreRes.isPassed {
		if existingProgress != nil && existingProgress.Status == domain.ProgressStatusCompleted {
			targetStatus = domain.ProgressStatusCompleted
		} else {
			targetStatus = domain.ProgressStatusInProgress
		}
	}

	err = s.txManager.WithTX(ctx, pgx.TxOptions{}, func(ctx context.Context) error {
		if err := s.saveLessonProgress(ctx, userID, lessonID, lesson.CourseID, progressScore, targetStatus); err != nil {
			return err
		}

		if s.quizRepo != nil {
			for _, es := range input.Essays {
				if err := s.quizRepo.SaveEssaySubmission(ctx, userID, lesson.CourseID, lessonID, es); err != nil {
					return fmt.Errorf("save essay submission: %w", err)
				}
			}
		}

		return s.recalculateCourseProgress(ctx, userID, lesson.CourseID)
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	validationResult.LessonID = lessonID
	validationResult.Status = string(targetStatus)
	validationResult.Score = scoreRes.finalScore
	validationResult.EarnedPoints = scoreRes.earnedPoints
	validationResult.TotalMaxPoints = scoreRes.totalPoints
	validationResult.IsPassed = scoreRes.isPassed

	return validationResult, nil
}

func (s *service) resolveLessonScore(ctx context.Context, userID, lessonID int64, lesson *domain.Lesson, input domain.CompleteLessonInput, validationResult *domain.LessonCompletionResult) (*scoreResolution, error) {
	hasPuckQuizzes := validationResult.TotalMaxPoints > 0
	isTestLessonType := lesson.Type == "test" || lesson.Type == "quiz"

	var dbQuiz *domain.Quiz
	if s.quizRepo != nil {
		if qz, qErr := s.quizRepo.GetQuizByLessonID(ctx, lessonID); qErr == nil {
			dbQuiz = qz
		}
	}

	isQuizLesson := isTestLessonType || hasPuckQuizzes || (input.AttemptID != nil) || (dbQuiz != nil)
	if !isQuizLesson {
		score := 100
		if input.Score != nil {
			score = *input.Score
		}
		return &scoreResolution{
			finalScore:   score,
			earnedPoints: score,
			totalPoints:  100,
			isPassed:     true,
		}, nil
	}

	switch {
	case input.AttemptID != nil:
		return s.resolveAttemptScore(ctx, userID, lessonID, *input.AttemptID, input.IsAbandoned)

	case hasPuckQuizzes:
		passingThreshold := 70
		if dbQuiz != nil && dbQuiz.PassingScor > 0 {
			passingThreshold = dbQuiz.PassingScor
		}
		return &scoreResolution{
			finalScore:   validationResult.Score,
			earnedPoints: validationResult.EarnedPoints,
			totalPoints:  validationResult.TotalMaxPoints,
			isPassed:     validationResult.Score >= passingThreshold,
		}, nil

	default:
		return &scoreResolution{
			finalScore:   0,
			earnedPoints: 0,
			totalPoints:  100,
			isPassed:     false,
		}, nil
	}
}

func (s *service) resolveAttemptScore(ctx context.Context, userID, lessonID int64, attemptID int64, isAbandoned bool) (*scoreResolution, error) {
	if s.quizRepo == nil {
		return nil, errorsAPP.ErrQuizNotFound
	}

	attempt, err := s.quizRepo.GetAttemptByID(ctx, attemptID)
	if err != nil {
		return nil, fmt.Errorf("get attempt: %w", err)
	}
	if attempt.UserID != userID {
		return nil, errorsAPP.ErrForbidden
	}

	quizObj, err := s.quizRepo.GetQuizByID(ctx, attempt.QuizID)
	if err != nil {
		return nil, fmt.Errorf("get quiz: %w", err)
	}
	if quizObj.LessonID != lessonID {
		return nil, fmt.Errorf("attempt belongs to different lesson: %w", errorsAPP.ErrForbidden)
	}

	if attempt.CompletedAt == nil || isAbandoned {
		now := time.Now()
		attempt.CompletedAt = &now
		curPts, sErr := s.quizRepo.SumAttemptPoints(ctx, attempt.ID)
		if sErr != nil {
			curPts = 0
		}
		totPts, tErr := s.quizRepo.GetQuizTotalPoints(ctx, quizObj.ID)
		if tErr != nil || totPts <= 0 {
			totPts = 1
		}
		attempt.Score = (curPts * 100) / totPts
		if attempt.Score > 100 {
			attempt.Score = 100
		}
		attempt.Passed = attempt.Score >= quizObj.PassingScor
		if uErr := s.quizRepo.UpdateAttempt(ctx, attempt); uErr != nil {
			return nil, fmt.Errorf("update attempt: %w", uErr)
		}
	}

	return &scoreResolution{
		finalScore:   attempt.Score,
		earnedPoints: attempt.Score,
		totalPoints:  100,
		isPassed:     attempt.Passed,
	}, nil
}

func (s *service) saveLessonProgress(ctx context.Context, userID, lessonID, courseID int64, progressScore int, targetStatus domain.ProgressStatus) error {
	err := s.progressRepo.UpdateLessonProgressStatus(ctx, userID, lessonID, targetStatus)
	if err != nil {
		if errors.Is(err, errorsAPP.ErrLessonProgressNotFound) {
			now := time.Now()
			lp := &domain.LessonProgress{
				UserID:      userID,
				LessonID:    lessonID,
				CourseID:    courseID,
				Status:      targetStatus,
				Score:       &progressScore,
				CompletedAt: &now,
				StartedAt:   &now,
				UpdatedAt:   now,
			}
			if cErr := s.progressRepo.CreateLessonProgress(ctx, lp); cErr != nil {
				return fmt.Errorf("create lesson progress: %w", cErr)
			}
			return nil
		}
		return fmt.Errorf("update lesson progress: %w", err)
	}

	if sErr := s.progressRepo.UpdateLessonProgressScore(ctx, userID, lessonID, progressScore); sErr != nil {
		return fmt.Errorf("update lesson progress score: %w", sErr)
	}
	return nil
}

func (s *service) recalculateCourseProgress(ctx context.Context, userID, courseID int64) error {
	courseLessons, err := s.lessonRepo.GetLessonsByCourseID(ctx, courseID)
	if err != nil {
		return fmt.Errorf("get course lessons: %w", err)
	}
	if len(courseLessons) == 0 {
		return nil
	}

	allProgress, pErr := s.progressRepo.GetAllLessonProgressByCourse(ctx, userID, courseID)
	if pErr != nil {
		return fmt.Errorf("get all lesson progress: %w", pErr)
	}

	completedCount := 0
	totalScore := 0
	scoreCount := 0
	for _, p := range allProgress {
		if p.Status == domain.ProgressStatusCompleted {
			completedCount++
			if p.Score != nil {
				totalScore += *p.Score
				scoreCount++
			}
		}
	}

	total := len(courseLessons)
	percent := (float64(completedCount) / float64(total)) * 100.0
	avgScore := 100.0
	if scoreCount > 0 {
		avgScore = float64(totalScore) / float64(scoreCount)
	}

	if err := s.progressRepo.UpsertCourseProgressWithScore(ctx, userID, courseID, completedCount, total, percent, avgScore); err != nil {
		return fmt.Errorf("upsert course progress: %w", err)
	}

	return nil
}
