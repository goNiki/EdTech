package repository

import (
	"context"
	"edtech/internal/domain"
	"edtech/internal/infrastructure/db"
	"time"
)

type UserRepository interface {
	CreateUser(ctx context.Context, q db.QueryExecutor, createUser domain.CreateUser) (*domain.User, error)
	GetUserByEmail(ctx context.Context, q db.QueryExecutor, email string) (*domain.User, error)
	GetUserByID(ctx context.Context, q db.QueryExecutor, userID int64) (*domain.User, error)
	GetUserByUserName(ctx context.Context, q db.QueryExecutor, username string) (*domain.User, error)
	ExistingByEmail(ctx context.Context, q db.QueryExecutor, email string) (bool, error)
	ExistingByUsernName(ctx context.Context, q db.QueryExecutor, username string) (bool, error)
	UpdateLastLogin(ctx context.Context, q db.QueryExecutor, userID int64, now time.Time) error
	UpdateProfile(ctx context.Context, q db.QueryExecutor, user *domain.User) error
	UpdatePassword(ctx context.Context, q db.QueryExecutor, userID int64, passHash string) error
	SetEmailVerified(ctx context.Context, q db.QueryExecutor, userID int64, verified bool) error
	UpdateRole(ctx context.Context, q db.QueryExecutor, userID int64, role domain.Role) error
	SetBannedStatus(ctx context.Context, q db.QueryExecutor, userID int64, isBanned bool) error
}

type CourseRepository interface {
	GetCourseBySlug(ctx context.Context, q db.QueryExecutor, slug string) (*domain.Course, error)
	CreateCourse(ctx context.Context, q db.QueryExecutor, course *domain.Course) (int64, error)
	GetCourseByID(ctx context.Context, q db.QueryExecutor, id int64) (*domain.Course, error)
	PublishCourse(ctx context.Context, q db.QueryExecutor, course *domain.Course) error
	UpdateCourse(ctx context.Context, q db.QueryExecutor, course *domain.Course) error
	ListPublicCourses(ctx context.Context, q db.QueryExecutor, pageSize int64, offset int64) ([]domain.Course, error)
	ListEnrolledCoursesAsStudent(ctx context.Context, q db.QueryExecutor, userID int64, pageSize int64, offset int64) ([]domain.Course, error)
	CountCourse(ctx context.Context, q db.QueryExecutor) (int, error)
	DeleteCourse(ctx context.Context, q db.QueryExecutor, courseID int64) error
}

type EnrolledRepository interface {
	UserExistCourse(ctx context.Context, q db.QueryExecutor, userID, courseid int64) (bool, error)
	EnrollUserToCourse(ctx context.Context, q db.QueryExecutor, enroll domain.EnrolledInCourse) error
	GetRoleUserInCourse(ctx context.Context, q db.QueryExecutor, userID, courceID int64) (string, error)
}

type LessonRepository interface {
	GetMaxPositionByCourseID(ctx context.Context, q db.QueryExecutor, courseId int64) (int64, error)
	CreateLesson(ctx context.Context, q db.QueryExecutor, lesson *domain.Lesson) error
	GetLessonByID(ctx context.Context, q db.QueryExecutor, id int64) (*domain.Lesson, error)
	UpdateLesson(ctx context.Context, q db.QueryExecutor, lesson *domain.Lesson) error
	DeleteLessonByID(ctx context.Context, q db.QueryExecutor, lessonID int64) error
	GetLessonsByCourseID(ctx context.Context, q db.QueryExecutor, courseID int64) ([]domain.Lesson, error)
}

type RefreshRepository interface {
	Save(ctx context.Context, q db.QueryExecutor, reftoken string, id int64, expiresAt time.Time) error
	GetByToken(ctx context.Context, q db.QueryExecutor, refreshToken string) (domain.RefreshTokenData, error)
	DeleteRefreshToken(ctx context.Context, q db.QueryExecutor, hashToken string) error
	DeleteAllByUserID(ctx context.Context, q db.QueryExecutor, userID int64) error
}

type PermissionsRepository interface {
	HasPermission(ctx context.Context, q db.QueryExecutor, roleName string, resource string, action string) (bool, error)
	GetRolePermissions(ctx context.Context, q db.QueryExecutor, rolename string) ([]string, error)
}
