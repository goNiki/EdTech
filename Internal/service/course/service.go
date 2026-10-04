package course

import (
	"edtech/internal/infrastructure/txmanager"
	"edtech/internal/repository"
	services "edtech/internal/service"
)

type service struct {
	courserepo    repository.CourseRepository
	sectionrepo   repository.SectionRepository
	lessonrepo    repository.LessonRepository
	accessService services.AccessService
	enrolledrepo  repository.EnrolledRepository
	txManager     txmanager.TransactionManager
}

func NewCourseService(courserepo repository.CourseRepository, sectionrepo repository.SectionRepository, lessonrepo repository.LessonRepository, accessService services.AccessService, enrolledrepo repository.EnrolledRepository, txManager txmanager.TransactionManager) *service {
	return &service{
		courserepo:    courserepo,
		sectionrepo:   sectionrepo,
		lessonrepo:    lessonrepo,
		accessService: accessService,
		enrolledrepo:  enrolledrepo,
		txManager:     txManager,
	}
}
