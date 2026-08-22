package service

import (
	"context"
	"edtech/internal/domain"
)

type CourseServices interface {
	CreateCourse(ctx context.Context, course *domain.Course) (int64, error)
	PublishCourse(ctx context.Context, userID int64, courseID int64) error
	GetCourseByID(ctx context.Context, id int64) (*domain.Course, error)
	GetCourseBySlug(ctx context.Context, slug string) (*domain.Course, error)
	ListPublicCourses(ctx context.Context, page int64, pageSize int64) (domain.PaginatedCourses, error)
	DeleteCourse(ctx context.Context, courseID int64, userID int64) error
}

type EnrolledServices interface {
	EnrollUserToCourse(context.Context, domain.EnrollUserRequest) error
	GetRoleUserInCource(ctx context.Context, userID, courseID int64) (string, error)
}

type LessonServices interface {
	CreateLesson(ctx context.Context, lesson *domain.Lesson) (int64, error)
	DeleteLesson(ctx context.Context, lessonID int64) error
	UpdateLesson(ctx context.Context, lesson *domain.Lesson) error
}

type AuthService interface {
	Register(ctx context.Context, input domain.RegisterInput) (*domain.User, error)
	Login(ctx context.Context, email, password string) (domain.AuthTokens, error)
	RefreshToken(ctx context.Context, refreshToken string) (domain.AuthTokens, error)
	Logout(ctx context.Context, refreshToken string) error
	GetCurrentUser(ctx context.Context, userID int64) (*domain.User, error)
	UpdateProfile(ctx context.Context, userID int64, input domain.UpdateProfileInput) (*domain.User, error)
	ChangePassword(ctx context.Context, userID int64, oldPassword, newPassword string) error
	VerifyEmail(ctx context.Context, userID int64) error
	ChangeUserRole(ctx context.Context, adminID int64, targetUserID int64, newRole domain.Role) error
	SetUserBanned(ctx context.Context, adminID int64, targetUserID int64, isBanned bool) error
}

type AccessService interface {
	CanAccessCourseObject(ctx context.Context, course *domain.Course, userID int64, action string) (bool, error)
	CanViewCourse(ctx context.Context, course *domain.Course, userID int64) (bool, error)
	CanEditCourse(ctx context.Context, course *domain.Course, userID int64) (bool, error)
	CanDeleteCourse(ctx context.Context, course *domain.Course, userID int64) (bool, error)
	CanPublishCourse(ctx context.Context, course *domain.Course, userID int64) (bool, error)
	CanSelfEnrollCourse(ctx context.Context, course *domain.Course, userID int64) (bool, error)
	CanManageCourseUsers(ctx context.Context, course *domain.Course, userID int64) (bool, error)
}
