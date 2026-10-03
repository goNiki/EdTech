package lesson

import (
	"edtech/internal/infrastructure/db"
	"edtech/internal/repository"
	"edtech/internal/service"
)

type lessonService struct {
	db            db.QueryExecutor
	courserepo    repository.CourseRepository
	lessonrepo    repository.LessonRepository
	sectionrepo   repository.SectionRepository
	accessService service.AccessService
}

func NewLessonService(
	courserepo repository.CourseRepository,
	lessonrepo repository.LessonRepository,
	sectionrepo repository.SectionRepository,
	accessService service.AccessService,
	database db.QueryExecutor,
) *lessonService {
	return &lessonService{
		db:            database,
		courserepo:    courserepo,
		lessonrepo:    lessonrepo,
		sectionrepo:   sectionrepo,
		accessService: accessService,
	}
}
