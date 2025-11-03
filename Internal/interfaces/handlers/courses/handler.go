package courses

import (
	"edtech/internal/infrastructure/validator"
	"edtech/internal/interfaces/middleware/auth"
	"edtech/internal/service"
)

type handler struct {
	courseService    service.CourseServices
	lessonService    service.LessonServices
	enrolmentService service.EnrolledServices
	validator        validator.Validator
	authMiddleware   auth.AuthMiddleware
}

func NewCourseHandler(courseService service.CourseServices, lessonService service.LessonServices, enrolmentService service.EnrolledServices, authMiddleware auth.AuthMiddleware) *handler {
	return &handler{
		courseService:    courseService,
		lessonService:    lessonService,
		enrolmentService: enrolmentService,
		validator:        *validator.NewValidator(),
		authMiddleware:   authMiddleware,
	}
}
