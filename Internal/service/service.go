package service

import (
	"context"
	"edtech/internal/domain"
	"edtech/internal/dto"
)

type CourseServices interface {
	CreateCourse(ctx context.Context, course *domain.Course) (int64, error)
	PublishCourse(ctx context.Context, id int64) error
}

type EnrolledServices interface {
	EnrollUserToCourse(context.Context, dto.EnrolleRequest) error
}

type LessonServices interface {
	CreateLesson(ctx context.Context, lesson *domain.Lesson) (int64, error)
	DeleteLesson(ctx context.Context, lessonID int64) error
	UpdateLesson(ctx context.Context, lesson *domain.Lesson) error
}

type AuthService interface {
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error)
	Login(ctx context.Context, email, password string) (*dto.LoginResponce, error)
}
