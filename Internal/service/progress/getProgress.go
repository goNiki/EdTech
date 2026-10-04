package progress

import (
	"context"
	"errors"
	"fmt"

	"edtech/internal/domain"
	errorsAPP "edtech/pkg/errors"

	"github.com/jackc/pgx/v5"
)

func (s *service) GetLessonProgress(ctx context.Context, userID int64, lessonID int64) (*domain.LessonProgress, error) {
	const op = "service.progress.GetLessonProgress"

	progress, err := s.progressRepo.GetLessonProgress(ctx, userID, lessonID)
	if err != nil {
		if errors.Is(err, errorsAPP.ErrLessonProgressNotFound) || errors.Is(err, pgx.ErrNoRows) {
			progress = &domain.LessonProgress{
				UserID:      userID,
				LessonID:    lessonID,
				Status:      domain.ProgressStatusNotStarted,
				Score:       nil,
				TimeSpent:   0,
				LastPos:     0,
				Submissions: make([]domain.LessonSubmissionDetail, 0),
			}
		} else {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	if progress.Submissions == nil {
		progress.Submissions = make([]domain.LessonSubmissionDetail, 0)
	}

	if s.quizRepo != nil {
		subs, sErr := s.quizRepo.GetLessonSubmissions(ctx, userID, lessonID)
		if sErr != nil {
			return nil, fmt.Errorf("%s: %w", op, sErr)
		}
		if subs != nil {
			progress.Submissions = subs
		}
	}

	return progress, nil
}

func (s *service) GetCourseProgress(ctx context.Context, userID int64, courseID int64) (*domain.CourseProgress, error) {
	const op = "service.progress.GetCourseProgress"

	progress, err := s.progressRepo.GetCourseProgress(ctx, userID, courseID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return progress, nil
}

func (s *service) GetAllLessonProgress(ctx context.Context, userID int64, courseID int64) ([]domain.LessonProgress, error) {
	const op = "service.progress.GetAllLessonProgress"

	progressList, err := s.progressRepo.GetAllLessonProgressByCourse(ctx, userID, courseID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return progressList, nil
}
