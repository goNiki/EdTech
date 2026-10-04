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

func (s *service) CompleteLesson(ctx context.Context, userID int64, lessonID int64, score *int, answers []domain.LessonAnswerSubmission, essays []domain.EssaySubmission) (*domain.LessonCompletionResult, error) {
	const op = "service.progress.CompleteLesson"

	lesson, err := s.lessonRepo.GetLessonByID(ctx, lessonID)
	if err != nil {
		return nil, fmt.Errorf("%s: get lesson: %w", op, err)
	}

	// Server-side validation of quizzes (anti-cheat verification)
	validationResult, vErr := quiz.ValidateQuizSubmission(lesson.Content, answers)
	if vErr != nil {
		return nil, fmt.Errorf("%s: validate quiz answers: %w", op, vErr)
	}

	finalScore := 100
	if validationResult.TotalMaxPoints > 0 {
		finalScore = validationResult.Score
	} else if score != nil {
		finalScore = *score
	}
	validationResult.Score = finalScore

	err = s.txManager.WithTX(ctx, pgx.TxOptions{}, func(ctx context.Context) error {
		// 1. Update or create lesson progress
		err = s.progressRepo.UpdateLessonProgressStatus(ctx, userID, lessonID, domain.ProgressStatusCompleted)
		if err != nil {
			if errors.Is(err, errorsAPP.ErrLessonProgressNotFound) {
				now := time.Now()
				lp := &domain.LessonProgress{
					UserID:      userID,
					LessonID:    lessonID,
					CourseID:    lesson.CourseID,
					Status:      domain.ProgressStatusCompleted,
					Score:       &finalScore,
					CompletedAt: &now,
					StartedAt:   &now,
					UpdatedAt:   now,
				}
				if cErr := s.progressRepo.CreateLessonProgress(ctx, lp); cErr != nil {
					return fmt.Errorf("create lesson progress: %w", cErr)
				}
			} else {
				return fmt.Errorf("update lesson progress: %w", err)
			}
		} else {
			if sErr := s.progressRepo.UpdateLessonProgressScore(ctx, userID, lessonID, finalScore); sErr != nil {
				return fmt.Errorf("update lesson progress score: %w", sErr)
			}
		}

		// 2. Handle Homework / Essay submissions
		for _, es := range essays {
			if err := s.quizRepo.SaveEssaySubmission(ctx, userID, lesson.CourseID, lessonID, es); err != nil {
				return fmt.Errorf("save essay submission: %w", err)
			}
		}

		// 3. Recalculate Course Progress
		courseLessons, err := s.lessonRepo.GetLessonsByCourseID(ctx, lesson.CourseID)
		if err != nil {
			return fmt.Errorf("get course lessons: %w", err)
		}

		if len(courseLessons) > 0 {
			allProgress, pErr := s.progressRepo.GetAllLessonProgressByCourse(ctx, userID, lesson.CourseID)
			if pErr != nil {
				return fmt.Errorf("get all lesson progress: %w", pErr)
			}

			completedCount := 0
			totalScore := 0
			scoreCount := 0
			for _, p := range allProgress {
				if p.Status == domain.ProgressStatusCompleted {
					completedCount++
					if p.Score != nil && *p.Score > 0 {
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

			if err := s.progressRepo.UpsertCourseProgressWithScore(ctx, userID, lesson.CourseID, completedCount, total, percent, avgScore); err != nil {
				return fmt.Errorf("upsert course progress: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return validationResult, nil
}

