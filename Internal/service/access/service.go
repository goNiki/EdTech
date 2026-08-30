package access

import (
	"edtech/internal/infrastructure/db"
	"edtech/internal/repository"
)

type service struct {
	db             db.QueryExecutor
	courserepo     repository.CourseRepository
	enrolledrepo   repository.EnrolledRepository
	permissionrepo repository.PermissionsRepository
}

func NewAccessService(courserepo repository.CourseRepository, enrolledrepo repository.EnrolledRepository, permissionrepo repository.PermissionsRepository, database db.QueryExecutor) *service {
	return &service{
		db:             database,
		courserepo:     courserepo,
		enrolledrepo:   enrolledrepo,
		permissionrepo: permissionrepo,
	}
}
