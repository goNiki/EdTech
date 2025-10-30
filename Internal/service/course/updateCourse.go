package course

import (
	"context"
	"edtech/internal/domain"
	"edtech/internal/infrastructure/logger"
	"edtech/internal/infrastructure/logger/sl"
	errorsAPP "edtech/pkg/errors"
	"edtech/pkg/utils"
	"errors"

	"github.com/jackc/pgx/v5"
)

func (s *service) UpdateCourse(ctx context.Context, course *domain.Course) error {
	const op = "usecase.course.updatecourse"

	log := logger.GetLogger(ctx, op)

	if err := utils.ValidateCourse(course.Title, course.Slug, int64(course.CreatedBy), course.Visibility, course.Status); err != nil {
		log.Error("error validate date", sl.Error(err))
		return err
	}

	if _, err := s.repo.GetCourseByID(ctx, int64(course.Id)); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error("course is not found", sl.Error(errorsAPP.ErrNotFoundCourse))
			return errorsAPP.ErrNotFoundCourse
		}
		log.Error("internal error BD:", sl.Error(err))
		return errorsAPP.ErrInternalDB
	}

	if err := s.repo.UpdateCourse(ctx, course); err != nil {
		if errors.Is(err, errorsAPP.ErrNothingToUpdate) {
			log.Error("nothing update", sl.Error(err))
			return errorsAPP.ErrNothingToUpdate
		}
		log.Error("internal error DB", sl.Error(err))
		return errorsAPP.ErrInternalDB
	}

	return nil

}
