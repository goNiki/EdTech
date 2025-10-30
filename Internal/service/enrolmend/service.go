package enrolmend

import (
	"edtech/internal/repository"
)

type service struct {
	courserepo   repository.CourseRepository
	enrolledrepo repository.EnrolledRepository
	userrepo     repository.UserRepository
}

func NewEnrolmentService(courseRepo *repository.CourseRepository, enrolledrepo *repository.EnrolledRepository, userrepo repository.UserRepository) *service {
	return &service{
		courserepo:   *courseRepo,
		enrolledrepo: *enrolledrepo,
		userrepo:     userrepo,
	}
}
