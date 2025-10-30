package course

import (
	"edtech/internal/repository"
)

type service struct {
	repo repository.CourseRepository
}

func NewCourseService(repo repository.CourseRepository) *service {
	return &service{
		repo: repo,
	}
}
