package lesson

import (
	"context"
	"edtech/internal/infrastructure/logger"
	"edtech/internal/infrastructure/logger/sl"
	errorsAPP "edtech/pkg/errors"
	"errors"
)

func (s *service) DeleteLesson(ctx context.Context, lessonID int64) error {

	const op = "usecase.course.lesson.deletelesson"

	log := logger.GetLogger(ctx, op)

	if err := s.lessonrepo.DeleteLessonByID(ctx, lessonID); err != nil {
		if errors.Is(err, errorsAPP.ErrNotFoundLesson) {
			log.Error("failed to delete lesson", sl.Error(err))
			return errorsAPP.ErrNotFoundLesson
		}
		log.Error("internal database err: ", sl.Error(err))
		return errorsAPP.ErrInternalDB
	}
	return nil
}
