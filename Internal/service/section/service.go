package section

import (
	"edtech/internal/infrastructure/txmanager"
	"edtech/internal/repository"
	services "edtech/internal/service"
)

var _ services.SectionServices = (*sectionService)(nil)

type sectionService struct {
	sectionRepo repository.SectionRepository
	lessonRepo  repository.LessonRepository
	txManager   txmanager.TransactionManager
}

func NewSectionService(
	sectionRepo repository.SectionRepository,
	lessonRepo repository.LessonRepository,
	txManager txmanager.TransactionManager,
) *sectionService {
	return &sectionService{
		sectionRepo: sectionRepo,
		lessonRepo:  lessonRepo,
		txManager:   txManager,
	}
}
