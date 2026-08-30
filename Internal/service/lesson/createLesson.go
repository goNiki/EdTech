package lesson

import (
	"context"
	"edtech/internal/domain"
	errorsAPP "edtech/pkg/errors"
	"edtech/pkg/utils"
	"errors"
	"fmt"
)

//ограничения:
//1. проверяем по id курса, есть ли вообще такой курс, +
//2. валидация обязательных полей +
//3. определение позиции (делать запрос в БД на количество уже действуюхих уроков в БД и выдача нового порядкового номера)
//4. вызываем функцию на запись данных в БД

func (s *service) CreateLesson(ctx context.Context, lesson *domain.Lesson) (int64, error) {

	const op = "service.lesson.CreateLesson"

	if err := utils.ValidateLesson(int64(lesson.CourseID), lesson.Title, lesson.Description); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	if _, err := s.courserepo.GetCourseByID(ctx, s.db, int64(lesson.CourseID)); err != nil {
		if errors.Is(err, errorsAPP.ErrNotFoundCourse) {
			return 0, fmt.Errorf("%s: %w", op, err)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	position, err := s.lessonrepo.GetMaxPositionByCourseID(ctx, s.db, int64(lesson.CourseID))
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	lesson.Position = position + 1

	err = s.lessonrepo.CreateLesson(ctx, s.db, lesson)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return int64(lesson.ID), nil
}
