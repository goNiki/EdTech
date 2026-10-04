package section

import (
	"edtech/internal/infrastructure/txmanager"
	"edtech/internal/repository"
	services "edtech/internal/service"
)

var _ services.SectionServices = (*sectionService)(nil)

type sectionService struct {
	sectionRepo   repository.SectionRepository
	lessonRepo    repository.LessonRepository
	courseRepo    repository.CourseRepository
	accessService services.AccessService
	txManager     txmanager.TransactionManager
}

func NewSectionService(
	sectionRepo repository.SectionRepository,
	lessonRepo repository.LessonRepository,
	courseRepo repository.CourseRepository,
	accessService services.AccessService,
	txManager txmanager.TransactionManager,
) *sectionService {
	return &sectionService{
		sectionRepo:   sectionRepo,
		lessonRepo:    lessonRepo,
		courseRepo:    courseRepo,
		accessService: accessService,
		txManager:     txManager,
	}
}
