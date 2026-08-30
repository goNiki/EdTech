package lesson

import (
	"edtech/internal/infrastructure/db"
	"edtech/internal/repository"
)

type service struct {
	db         db.QueryExecutor
	courserepo repository.CourseRepository
	lessonrepo repository.LessonRepository
}

func NewLessonService(courserepo repository.CourseRepository, lessonrepo repository.LessonRepository, database db.QueryExecutor) *service {
	return &service{
		db:         database,
		courserepo: courserepo,
		lessonrepo: lessonrepo,
	}
}
