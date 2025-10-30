package course

import (
	"context"
	"edtech/internal/infrastructure/logger"
	"edtech/internal/infrastructure/logger/sl"
	errorsAPP "edtech/pkg/errors"
	"errors"

	"github.com/jackc/pgx/v5"
)

func (s *service) PublishCourse(ctx context.Context, id int64) error {

	const op = "usecase.course.publishcourse"

	log := logger.GetLogger(ctx, op)

	course, err := s.repo.GetCourseByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error("course is not found", sl.Error(err))
			return errorsAPP.ErrNotFoundCourse
		}
		log.Error("Error DB", sl.Error(err))
		return errorsAPP.ErrInternalDB
	}
	//TODO еще дополнительную проверку, чтобы был хотя бы 1 урок на курсе, чтобы его можно было опубликовать
	err = course.Publish()
	if err != nil {
		log.Error("", sl.Error(err))
		return errorsAPP.ErrCourseAlredyPublished
	}

	if err := s.repo.PublishCourse(ctx, course); err != nil {
		log.Error("Error DB", sl.Error(err))
		return errorsAPP.ErrInternalDB
	}

	return nil

}
