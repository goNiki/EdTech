package service

import (
	"context"
	"edtech/internal/domain"
)

type CourseServices interface {
	CreateCourse(ctx context.Context, course *domain.Course) (*domain.Course, error)
	PublishCourse(ctx context.Context, userID int64, courseID int64) error
	ArchiveCourse(ctx context.Context, userID int64, courseID int64) error
	GetCourseByID(ctx context.Context, id int64) (*domain.Course, error)
	GetCourseBySlug(ctx context.Context, slug string) (*domain.Course, error)
	ListPublicCourses(ctx context.Context, pagination domain.Pagination, filter domain.CourseFilter) (domain.PaginatedCourses, error)
	ListMyCourses(ctx context.Context, input *domain.InputListMyCourse) (domain.PaginatedCourses, error)
	GetCourseStructure(ctx context.Context, courseID int64) (domain.CourseStructure, error)
	UpdateCourse(ctx context.Context, courseID int64, userID int64, input domain.UpdateCourseInput) (*domain.Course, error)
	DeleteCourse(ctx context.Context, courseID int64, userID int64) error
}

type EnrolledServices interface {
	SelfEnrollCourse(context.Context, domain.SelfEnrollRequest) error
	TeacherEnrollCourse(context.Context, domain.TeacherEnrollRequest) error
	GetRoleUserInCource(ctx context.Context, userID, courseID int64) (string, error)
	UnenrollUser(ctx context.Context, userID int64, courseID int64) error
	ListCourseStudents(ctx context.Context, courseID int64, page int64, pageSize int64) ([]domain.User, int, error)
	ChangeUserRole(ctx context.Context, courseID int64, targetUserID int64, newRole string) error
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
	BuildCoursePermissions(ctx context.Context, course *domain.Course, userID int64) (domain.CoursePermissions, error)
	CanAccessCourseObject(ctx context.Context, course *domain.Course, userID int64, action string) (bool, error)
	CanViewCourse(ctx context.Context, course *domain.Course, userID int64) (bool, error)
	CanEditCourse(ctx context.Context, course *domain.Course, userID int64) (bool, error)
	CanDeleteCourse(ctx context.Context, course *domain.Course, userID int64) (bool, error)
	CanPublishCourse(ctx context.Context, course *domain.Course, userID int64) (bool, error)
	CanSelfEnrollCourse(ctx context.Context, course *domain.Course, userID int64) (bool, error)
	CanManageCourseUsers(ctx context.Context, course *domain.Course, userID int64) (bool, error)
}

type ProgressServices interface {
	StartLesson(ctx context.Context, userID int64, lessonID int64) error
	UpdateLessonProgress(ctx context.Context, userID int64, lessonID int64, input domain.UpdateProgressInput) error
	CompleteLesson(ctx context.Context, userID int64, lessonID int64) error
	GetLessonProgress(ctx context.Context, userID int64, lessonID int64) (*domain.LessonProgress, error)
	GetCourseProgress(ctx context.Context, userID int64, courseID int64) (*domain.CourseProgress, error)
	GetAllLessonProgress(ctx context.Context, userID int64, courseID int64) ([]domain.LessonProgress, error)
}

type QuizServices interface {
	CreateQuiz(ctx context.Context, userID int64, quiz *domain.Quiz) (*domain.Quiz, error)
	StartAttempt(ctx context.Context, userID int64, quizID int64) (*domain.QuizAttempt, error)
	SubmitAttempt(ctx context.Context, userID int64, attemptID int64, answers []domain.QuizAttemptAnswer) (*domain.QuizAttempt, error)
	ListAttemptsForGrading(ctx context.Context, userID int64, courseID int64, quizID *int64, page int, pageSize int) ([]domain.QuizAttempt, int64, error)
	GradeAttemptAnswer(ctx context.Context, teacherID int64, attemptID int64, answerID int64, points int, feedback *string) (*domain.QuizAttempt, error)
}
