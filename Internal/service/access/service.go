package access

import (
	"edtech/internal/repository"
)

type service struct {
	courserepo     repository.CourseRepository
	enrolledrepo   repository.EnrolledRepository
	permissionrepo repository.PermissionsRepository
}

func NewAccessService(courserepo repository.CourseRepository, enrolledrepo repository.EnrolledRepository, permissionrepo repository.PermissionsRepository) *service {
	return &service{
		courserepo:     courserepo,
		enrolledrepo:   enrolledrepo,
		permissionrepo: permissionrepo,
	}
}
