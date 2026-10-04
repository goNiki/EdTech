package lesson

import (
	"edtech/internal/repository"
	"edtech/internal/service"
)

type lessonService struct {
	courserepo    repository.CourseRepository
	lessonrepo    repository.LessonRepository
	sectionrepo   repository.SectionRepository
	accessService service.AccessService
	progressrepo  repository.ProgressRepository
}

func NewLessonService(
	courserepo repository.CourseRepository,
	lessonrepo repository.LessonRepository,
	sectionrepo repository.SectionRepository,
	accessService service.AccessService,
	progressrepo repository.ProgressRepository,
) *lessonService {
	return &lessonService{
		courserepo:    courserepo,
		lessonrepo:    lessonrepo,
		sectionrepo:   sectionrepo,
		accessService: accessService,
		progressrepo:  progressrepo,
	}
}

