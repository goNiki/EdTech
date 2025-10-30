package course

import (
	"context"
	"edtech/internal/domain"
	"edtech/internal/infrastructure/logger"
	"edtech/internal/infrastructure/logger/sl"
	errorsAPP "edtech/pkg/errors"
	"edtech/pkg/utils"
	"errors"
)

func (s *service) CreateCourse(ctx context.Context, course *domain.Course) (int64, error) {

	const op = "usecase.course_uc.createcourse"

	log := logger.GetLogger(ctx, op)

	if err := utils.ValidateCourse(course.Title, course.Slug, int64(course.CreatedBy), course.Visibility, course.Status); err != nil {
		log.Error("error validate date", sl.Error(err))
		return 0, err
	}

	if _, err := s.repo.GetCourseBySlug(ctx, course.Slug); !errors.Is(err, errorsAPP.ErrNotFoundCourse) {
		if err == nil {
			log.Error("slug already exists:", sl.Error(errorsAPP.ErrSlugAlreadyExists))
			return 0, errorsAPP.ErrSlugAlreadyExists
		}
		log.Error("error DB", sl.Error(err))
		return 0, errorsAPP.ErrInternalDB
	}

	course.Status = "draft"

	courseID, err := s.repo.CreateCourse(ctx, course)
	if err != nil {
		log.Error("database error", sl.Error(err))
		return 0, errorsAPP.ErrInternalDB
	}

	return courseID, nil

}
