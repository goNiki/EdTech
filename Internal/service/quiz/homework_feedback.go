package quiz

import (
	"context"
	"fmt"

	"edtech/internal/domain"
	errorsAPP "edtech/pkg/errors"
)

func (s *service) GetStudentHomeworkFeedback(ctx context.Context, userID, lessonID int64) (*domain.StudentHomeworkFeedback, error) {
	const op = "service.quiz.GetStudentHomeworkFeedback"

	if userID <= 0 {
		return nil, fmt.Errorf("%s: %w", op, errorsAPP.ErrUnauthorized)
	}

	if lessonID <= 0 {
		return nil, fmt.Errorf("%s: %w", op, errorsAPP.ErrInvalidURLParam)
	}

	// 1. Проверяем существование урока
	lesson, err := s.lessonRepo.GetLessonByID(ctx, lessonID)
	if err != nil {
		return nil, fmt.Errorf("%s: get lesson: %w", op, err)
	}

	// 2. Проверяем существование курса и доступ пользователя к курсу
	course, err := s.courseRepo.GetCourseByID(ctx, lesson.CourseID)
	if err != nil {
		return nil, fmt.Errorf("%s: get course: %w", op, err)
	}

	if s.accessService != nil {
		canView, aErr := s.accessService.CanViewCourse(ctx, course, userID)
		if aErr != nil {
			return nil, fmt.Errorf("%s: check course access: %w", op, aErr)
		}
		if !canView {
			return nil, fmt.Errorf("%s: %w", op, errorsAPP.ErrForbidden)
		}
	}

	// 3. Получаем рецензию и детали сдачи из репозитория
	feedback, err := s.quizRepo.GetStudentHomeworkFeedback(ctx, userID, lessonID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return feedback, nil
}
