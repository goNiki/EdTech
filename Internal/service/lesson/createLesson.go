package lesson

import (
	"context"
	"edtech/internal/domain"
	"edtech/internal/infrastructure/logger"
	"edtech/internal/infrastructure/logger/sl"
	errorsAPP "edtech/pkg/errors"
	"edtech/pkg/utils"
	"errors"
)

//ограничения:
//1. проверяем по id курса, есть ли вообще такой курс, +
//2. валидация обязательных полей +
//3. определение позиции (делать запрос в БД на количество уже действуюхих уроков в БД и выдача нового порядкового номера)
//4. вызываем функцию на запись данных в БД

func (s *service) CreateLesson(ctx context.Context, lesson *domain.Lesson) (int64, error) {

	const op = "usecase.course.createlesson"

	log := logger.GetLogger(ctx, op)

	if err := utils.ValidateLesson(int64(lesson.CourseID), lesson.Title, lesson.Description); err != nil {
		log.Error("validate Error", sl.Error(err))
		return 0, err
	}

	if _, err := s.courserepo.GetCourseByID(ctx, int64(lesson.CourseID)); err != nil {
		if errors.Is(err, errorsAPP.ErrNotFoundCourse) {
			log.Error("%w", sl.Error(err))
			return 0, errorsAPP.ErrNotFoundCourse
		}
	}

	position, err := s.lessonrepo.GetMaxPositionByCourseID(ctx, int64(lesson.CourseID))
	if err != nil {
		log.Error("internal database error: ", sl.Error(err))
		return 0, errorsAPP.ErrInternalDB
	}

	lesson.Position = int(position) + 1

	err = s.lessonrepo.CreateLesson(ctx, lesson)
	if err != nil {
		log.Error("internal database error", sl.Error(err))
		return 0, errorsAPP.ErrInternalDB
	}

	return int64(lesson.ID), nil

}
