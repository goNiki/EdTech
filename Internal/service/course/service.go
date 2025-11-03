package course

import (
	"edtech/internal/repository"
	services "edtech/internal/service"
)

type service struct {
	repo          repository.CourseRepository
	accessService services.AccessService
}

func NewCourseService(repo repository.CourseRepository, accessService services.AccessService) *service {
	return &service{
		repo:          repo,
		accessService: accessService,
	}
}
