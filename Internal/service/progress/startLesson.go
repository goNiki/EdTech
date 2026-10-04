package progress

import (
	"context"
	"fmt"
	"time"

	"edtech/internal/domain"

	"github.com/jackc/pgx/v5"
)

func (s *service) StartLesson(ctx context.Context, userID int64, lessonID int64) error {
	const op = "service.progress.StartLesson"

	lesson, err := s.lessonRepo.GetLessonByID(ctx, lessonID)
	if err != nil {
		return fmt.Errorf("%s: get lesson: %w", op, err)
	}

	lessons, err := s.lessonRepo.GetLessonsByCourseID(ctx, lesson.CourseID)
	if err != nil {
		return fmt.Errorf("%s: get lessons by course: %w", op, err)
	}

	totalLessons := len(lessons)
	if totalLessons == 0 {
		totalLessons = 1
	}

	now := time.Now()

	err = s.txManager.WithTX(ctx, pgx.TxOptions{}, func(ctx context.Context) error {
		courseProg := &domain.CourseProgress{
			UserID:         userID,
			CourseID:       lesson.CourseID,
			TotalLessons:   totalLessons,
			CompletedLess:  0,
			Percent:        0,
			TotalWatchTime: 0,
			StartedAt:      &now,
			LastAccessedAt: now,
		}
		if err := s.progressRepo.CreateCourseProgress(ctx, courseProg); err != nil {
			return err
		}

		lessonProg := &domain.LessonProgress{
			UserID:    userID,
			LessonID:  lessonID,
			CourseID:  lesson.CourseID,
			Status:    domain.ProgressStatusInProgress,
			StartedAt: &now,
		}
		if err := s.progressRepo.CreateLessonProgress(ctx, lessonProg); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
