package service

import (
	"context"
	"edtech/internal/domain"
	"edtech/internal/dto"
)

type CourseServices interface {
	CreateCourse(ctx context.Context, course *domain.Course) (int64, error)
	PublishCourse(ctx context.Context, id int64) error
	GetCourseByID(ctx context.Context, id int64) (*domain.Course, error)
	GetCourseBySlug(ctx context.Context, slug string) (*domain.Course, error)
	ListCourses(ctx context.Context, page int, pageSize int) (dto.PaginatedCourses, error)
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

type AccessService interface {
	// Универсальный метод для проверки доступа
	CanAccessCourseObject(ctx context.Context, course *domain.Course, userID int64, action string) (bool, error)

	// Специфичные методы для разных действий
	CanViewCourse(ctx context.Context, course *domain.Course, userID int64) (bool, error)
	CanEditCourse(ctx context.Context, course *domain.Course, userID int64) (bool, error)
	CanDeleteCourse(ctx context.Context, course *domain.Course, userID int64) (bool, error)
	CanPublishCourse(ctx context.Context, course *domain.Course, userID int64) (bool, error)

	// Самостоятельная запись студента на курс
	CanSelfEnrollCourse(ctx context.Context, course *domain.Course, userID int64) (bool, error)

	// Управление пользователями курса (добавление/удаление учителем)
	CanManageCourseUsers(ctx context.Context, course *domain.Course, userID int64) (bool, error)
}
