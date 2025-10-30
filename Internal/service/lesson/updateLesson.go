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

//порядоку дейстий:
//проверяем наличие урока по ID, достаем исходые данные урока.
//валидируем значение title, если значение пустое, то мы выдаем ошибку что поле пустое(есть вариант сделать чтобы подставлялось в таких случаях предыдущее значение title), полученные выше.
//content меняется в любом случае
//вызываем функцию изменения в репозитории, меняем 2 значения title (даже если значения одинаковые) и content

func (s *service) UpdateLesson(ctx context.Context, lesson *domain.Lesson) error {
	const op = "usecase.course.lesson.updatelesson"

	log := logger.GetLogger(ctx, op)

	if _, err := s.lessonrepo.GetLessonByID(ctx, lesson.ID); err != nil {
		if errors.Is(err, errorsAPP.ErrNotFoundLesson) {
			log.Error("", sl.Error(err))
			return errorsAPP.ErrNotFoundLesson
		}
		log.Error("internal database err: ", sl.Error(err))
		return errorsAPP.ErrInternalDB
	}
	//костыль- передаю ид урока, вместо ид курса.
	if err := utils.ValidateLesson(int64(lesson.ID), lesson.Title, lesson.Description); err != nil {
		log.Error("", sl.Error(err))
		return err
	}

	err := s.lessonrepo.UpdateLesson(ctx, lesson)
	if err != nil {
		if err == errorsAPP.ErrNothingToUpdate {
			log.Error("nothing to update for lesson")
			return errorsAPP.ErrNothingToUpdate
		}
		log.Error("Internal database error", sl.Error(err))
		return errorsAPP.ErrInternalDB
	}

	return nil
}
