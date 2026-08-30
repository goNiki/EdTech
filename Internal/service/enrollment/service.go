package enrollment

import (
	"edtech/internal/infrastructure/db"
	"edtech/internal/infrastructure/txmanager"
	"edtech/internal/repository"
	services "edtech/internal/service"
)

type service struct {
	db            db.QueryExecutor
	txManager     txmanager.TransactionManager
	courserepo    repository.CourseRepository
	enrolledrepo  repository.EnrolledRepository
	userrepo      repository.UserRepository
	accessService services.AccessService
}

func NewEnrolmentService(courseRepo repository.CourseRepository, enrolledrepo repository.EnrolledRepository, userrepo repository.UserRepository, accessService services.AccessService, database db.QueryExecutor, txManager txmanager.TransactionManager) *service {
	return &service{
		db:            database,
		txManager:     txManager,
		courserepo:    courseRepo,
		enrolledrepo:  enrolledrepo,
		userrepo:      userrepo,
		accessService: accessService,
	}
}
