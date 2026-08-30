package quiz

import (
	"context"
	"fmt"

	"edtech/internal/domain"
	errorsAPP "edtech/pkg/errors"
)

func (s *service) ListAttemptsForGrading(ctx context.Context, userID int64, courseID int64, quizID *int64, page int, pageSize int) ([]domain.QuizAttempt, int64, error) {
	const op = "service.quiz.ListAttemptsForGrading"

	course, err := s.courseRepo.GetCourseByID(ctx, s.db, courseID)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: get course: %w", op, err)
	}

	canEdit, err := s.accessService.CanEditCourse(ctx, course, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrCheckingPermissions, err)
	}
	if !canEdit {
		return nil, 0, fmt.Errorf("%s: %w", op, errorsAPP.ErrForbidden)
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	attempts, total, err := s.quizRepo.ListAttemptsForGrading(ctx, s.db, courseID, quizID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	return attempts, total, nil
}
