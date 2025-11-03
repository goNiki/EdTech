package repository

import (
	"context"
	"edtech/internal/domain"
	"time"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
}

type CourseRepository interface {
	GetCourseBySlug(ctx context.Context, slug string) (*domain.Course, error)
	CreateCourse(ctx context.Context, course *domain.Course) (int64, error)
	GetCourseByID(ctx context.Context, id int64) (*domain.Course, error)
	PublishCourse(ctx context.Context, course *domain.Course) error
	UpdateCourse(ctx context.Context, course *domain.Course) error
	ListCourses(ctx context.Context, pageSize int, offset int) ([]domain.Course, error)
	CountCourse(ctx context.Context) (int, error)
}

type EnrolledRepository interface {
	UserExistCourse(ctx context.Context, userID, courseid int64) (bool, error)
	EnrollUserToCourse(ctx context.Context, enroll domain.EnrolledInCourse) error
	GetRoleUserInCourse(ctx context.Context, userID, courceID int64) (string, error)
}

type LessonRepository interface {
	GetMaxPositionByCourseID(ctx context.Context, courseId int64) (int64, error)
	CreateLesson(ctx context.Context, lesson *domain.Lesson) error
	GetLessonByID(ctx context.Context, id int) (*domain.Lesson, error)
	UpdateLesson(ctx context.Context, lesson *domain.Lesson) error
	DeleteLessonByID(ctx context.Context, lessonID int64) error
	GetLessonsByCourseID(ctx context.Context, courseID int) ([]domain.Lesson, error)
}

type RefreshRepository interface {
	Save(ctx context.Context, reftoken string, id int64, expiresAt time.Time) error
}

type PermissionsRepository interface {
	HasPermission(ctx context.Context, roleName string, resource string, action string) (bool, error)
	GetRolePermissions(ctx context.Context, rolename string) ([]string, error)
}
