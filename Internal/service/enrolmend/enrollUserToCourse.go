package enrolmend

import (
	"context"
	"edtech/internal/domain"
	"edtech/internal/dto"
	"edtech/internal/infrastructure/logger"
	"edtech/internal/infrastructure/logger/sl"
	errorsAPP "edtech/pkg/errors"
	"edtech/pkg/utils"
	"errors"
)

type EnrolledServices interface {
}

func (s *service) EnrollUserToCourse(ctx context.Context, enrol dto.EnrolleRequest) error {
	const op = "usecase.enrolled.enrollusertocourse"

	log := logger.GetLogger(ctx, op)

	if err := utils.ValidateEnrolle(enrol); err != nil {
		log.Error("Error: ", sl.Error(err))
		return errorsAPP.ErrFailEnroleValidate
	}

	user, err := s.userrepo.GetUserByEmail(ctx, enrol.UserEmail)
	if err != nil {
		if errors.Is(err, errorsAPP.ErrUserNotFound) {
			log.Error("user not found", sl.Error(err))
			return errorsAPP.ErrUserNotFound
		}
		log.Error("nternal error", sl.Error(err))
		return errorsAPP.ErrInternalDB
	}

	if _, err := s.courserepo.GetCourseByID(ctx, int64(enrol.CourseID)); err != nil {
		if errors.Is(err, errorsAPP.ErrNotFoundCourse) {
			log.Error("course not found", sl.Error(err))
			return errorsAPP.ErrNotFoundCourse
		}
		log.Error("internal error", sl.Error(err))
		return errorsAPP.ErrInternalDB
	}

	exist, err := s.enrolledrepo.UserExistCourse(ctx, int(user.ID), enrol.CourseID)
	if err != nil {
		log.Error("internal error: ", sl.Error(err))
		return errorsAPP.ErrInternalDB
	}
	if exist {
		log.Error("user already enrolled: ", sl.Error(err))
		return errorsAPP.ErrUserAlreadyEnrolled
	}

	enroll := domain.EnrolledInCourse{
		UserID:   int(user.ID),
		CourseID: enrol.CourseID,
		Role:     enrol.Role,
	}

	if err := s.enrolledrepo.EnrollUserToCourse(ctx, enroll); err != nil {
		log.Error("internal error", sl.Error(err))
		return errorsAPP.ErrInternalDB
	}

	return nil
}
