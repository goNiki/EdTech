package course

import (
	"edtech/internal/repository"
	services "edtech/internal/service"
)

type service struct {
	courserepo    repository.CourseRepository
	accessService services.AccessService
}

func NewCourseService(courserepo repository.CourseRepository, accessService services.AccessService) *service {
	return &service{
		courserepo:    courserepo,
		accessService: accessService,
	}
}
