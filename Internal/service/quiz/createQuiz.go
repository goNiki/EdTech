package quiz

import (
	"context"
	"fmt"

	"edtech/internal/domain"
	errorsAPP "edtech/pkg/errors"
)

func (s *service) CreateQuiz(ctx context.Context, userID int64, quiz *domain.Quiz) (*domain.Quiz, error) {
	const op = "service.quiz.CreateQuiz"

	if quiz == nil {
		return nil, fmt.Errorf("%s: %w", op, errorsAPP.ErrQuizValidation)
	}

	if err := quiz.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Проверяем существование урока
	lesson, err := s.lessonRepo.GetLessonByID(ctx, quiz.LessonID)
	if err != nil {
		return nil, fmt.Errorf("%s: get lesson: %w", op, err)
	}

	// Проверяем права на редактирование курса
	course, err := s.courseRepo.GetCourseByID(ctx, lesson.CourseID)
	if err != nil {
		return nil, fmt.Errorf("%s: get course: %w", op, err)
	}

	canEdit, err := s.accessService.CanEditCourse(ctx, course, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrCheckingPermissions, err)
	}
	if !canEdit {
		return nil, fmt.Errorf("%s: %w", op, errorsAPP.ErrForbidden)
	}

	createdQuiz, err := s.quizRepo.CreateQuiz(ctx, quiz)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return createdQuiz, nil
}
