package progress

import (
	"context"
	"fmt"

	"edtech/internal/domain"
)

func (s *service) GetLessonProgress(ctx context.Context, userID int64, lessonID int64) (*domain.LessonProgress, error) {
	const op = "service.progress.GetLessonProgress"

	progress, err := s.progressRepo.GetLessonProgress(ctx, s.db, userID, lessonID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return progress, nil
}

func (s *service) GetCourseProgress(ctx context.Context, userID int64, courseID int64) (*domain.CourseProgress, error) {
	const op = "service.progress.GetCourseProgress"

	progress, err := s.progressRepo.GetCourseProgress(ctx, s.db, userID, courseID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return progress, nil
}

func (s *service) GetAllLessonProgress(ctx context.Context, userID int64, courseID int64) ([]domain.LessonProgress, error) {
	const op = "service.progress.GetAllLessonProgress"

	progressList, err := s.progressRepo.GetAllLessonProgressByCourse(ctx, s.db, userID, courseID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return progressList, nil
}
