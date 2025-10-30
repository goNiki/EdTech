package lesson

import (
	"edtech/internal/repository"
)

type service struct {
	courserepo repository.CourseRepository
	lessonrepo repository.LessonRepository
}

func NewLessonService(courserepo repository.CourseRepository, lessonrepo repository.LessonRepository) *service {
	return &service{
		courserepo: courserepo,
		lessonrepo: lessonrepo,
	}
}
