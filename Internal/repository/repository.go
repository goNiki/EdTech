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
	CreateCourse(ctx context.Context, q db.QueryExecutor, course *domain.Course) (*domain.Course, error)
	GetCourseByID(ctx context.Context, q db.QueryExecutor, id int64) (*domain.Course, error)
	ArchiveCourse(ctx context.Context, q db.QueryExecutor, course *domain.Course) error
	PublishCourse(ctx context.Context, q db.QueryExecutor, course *domain.Course) error
	UpdateCourseStatus(ctx context.Context, q db.QueryExecutor, courseID int64, status string) error
	UpdateCourse(ctx context.Context, q db.QueryExecutor, course *domain.Course) error
	ListPublicCourses(ctx context.Context, q db.QueryExecutor, pagination domain.Pagination, filter domain.CourseFilter) ([]domain.Course, error)
	ListEnrolledCourses(ctx context.Context, q db.QueryExecutor, input *domain.InputListMyCourse) ([]domain.Course, error)
	CountEnrolledCourses(ctx context.Context, q db.QueryExecutor, input *domain.InputListMyCourse) (int64, error)
	CountCourses(ctx context.Context, q db.QueryExecutor, filter domain.CourseFilter) (int64, error)
	DeleteCourse(ctx context.Context, q db.QueryExecutor, courseID int64) error
	ExistingBySlug(ctx context.Context, q db.QueryExecutor, slug string) (bool, error)
	ExistsBySlug(ctx context.Context, q db.QueryExecutor, slug string) (bool, error)
	IncrementEnrolledCount(ctx context.Context, q db.QueryExecutor, courseID int64) error
	DecrementEnrolledCount(ctx context.Context, q db.QueryExecutor, courseID int64) error
}

type EnrolledRepository interface {
	UserExistCourse(ctx context.Context, q db.QueryExecutor, userID, courseid int64) (bool, error)
	EnrollUserToCourse(ctx context.Context, q db.QueryExecutor, enroll domain.EnrolledInCourse) error
	GetRoleUserInCourse(ctx context.Context, q db.QueryExecutor, userID, courceID int64) (string, error)
	UnenrollUser(ctx context.Context, q db.QueryExecutor, userID, courseID int64) error
	ListCourseStudents(ctx context.Context, q db.QueryExecutor, courseID int64, limit, offset int64) ([]domain.User, int, error)
	ListCourseStudentsWithProgress(ctx context.Context, q db.QueryExecutor, courseID int64, limit, offset int64) ([]domain.CourseStudentItem, int64, error)
	ChangeUserRole(ctx context.Context, q db.QueryExecutor, courseID int64, targetUserID int64, newRole string) error
}

type SectionRepository interface {
	CreateSection(ctx context.Context, q db.QueryExecutor, section *domain.Section) (*domain.Section, error)
	GetSectionByID(ctx context.Context, q db.QueryExecutor, id int64) (*domain.Section, error)
	ListSectionsByCourseID(ctx context.Context, q db.QueryExecutor, courseID int64) ([]domain.Section, error)
	UpdateSection(ctx context.Context, q db.QueryExecutor, section *domain.Section) error
	UpdateSectionStatus(ctx context.Context, q db.QueryExecutor, sectionID int64, status string) error
	UpdateStatusByCourseID(ctx context.Context, q db.QueryExecutor, courseID int64, status string) error
	ReorderSections(ctx context.Context, q db.QueryExecutor, courseID int64, sectionIDs []int64) error
	DeleteSection(ctx context.Context, q db.QueryExecutor, sectionID int64) error
	GetMaxPositionByCourseID(ctx context.Context, q db.QueryExecutor, courseID int64) (int, error)
}

type LessonRepository interface {
	GetMaxPositionByCourseID(ctx context.Context, q db.QueryExecutor, courseId int64) (int64, error)
	CreateLesson(ctx context.Context, q db.QueryExecutor, lesson *domain.Lesson) error
	GetLessonByID(ctx context.Context, q db.QueryExecutor, id int64) (*domain.Lesson, error)
	UpdateLesson(ctx context.Context, q db.QueryExecutor, lesson *domain.Lesson) error
	UpdateLessonStatus(ctx context.Context, q db.QueryExecutor, lessonID int64, status string) error
	UpdateStatusByCourseID(ctx context.Context, q db.QueryExecutor, courseID int64, status string) error
	UpdateStatusBySectionID(ctx context.Context, q db.QueryExecutor, sectionID int64, status string) error
	ReorderLessons(ctx context.Context, q db.QueryExecutor, sectionID *int64, lessonIDs []int64) error
	DeleteLessonByID(ctx context.Context, q db.QueryExecutor, lessonID int64) error
	GetLessonsByCourseID(ctx context.Context, q db.QueryExecutor, courseID int64) ([]domain.Lesson, error)
}

type AnalyticsRepository interface {
	GetCourseAnalyticsSummary(ctx context.Context, q db.QueryExecutor, courseID int64) (domain.CourseAnalyticsSummary, error)
	ListPendingHomeworks(ctx context.Context, q db.QueryExecutor, courseID int64, limit, offset int64) ([]domain.PendingHomeworkItem, int64, error)
	GetStudentDrilldown(ctx context.Context, q db.QueryExecutor, userID, courseID int64) (*domain.StudentDrilldownReport, error)
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

type ProgressRepository interface {
	CreateLessonProgress(ctx context.Context, q db.QueryExecutor, progress *domain.LessonProgress) error
	GetLessonProgress(ctx context.Context, q db.QueryExecutor, userID, lessonID int64) (*domain.LessonProgress, error)
	UpdateLessonProgressTime(ctx context.Context, q db.QueryExecutor, userID, lessonID int64, additionalTime int, lastPos int) error
	UpdateLessonProgressStatus(ctx context.Context, q db.QueryExecutor, userID, lessonID int64, status domain.ProgressStatus) error
	UpdateLessonProgressScore(ctx context.Context, q db.QueryExecutor, userID, lessonID int64, score int) error

	CreateCourseProgress(ctx context.Context, q db.QueryExecutor, progress *domain.CourseProgress) error
	UpsertCourseProgress(ctx context.Context, q db.QueryExecutor, userID, courseID int64, completedLessons, totalLessons int, percentage float64) error
	UpsertCourseProgressWithScore(ctx context.Context, q db.QueryExecutor, userID, courseID int64, completedLessons, totalLessons int, percentage float64, averageScore float64) error
	GetCourseProgress(ctx context.Context, q db.QueryExecutor, userID, courseID int64) (*domain.CourseProgress, error)
	GetAllLessonProgressByCourse(ctx context.Context, q db.QueryExecutor, userID, courseID int64) ([]domain.LessonProgress, error)
}

type QuizRepository interface {
	CreateQuiz(ctx context.Context, q db.QueryExecutor, quiz *domain.Quiz) (*domain.Quiz, error)
	GetQuizByID(ctx context.Context, q db.QueryExecutor, id int64) (*domain.Quiz, error)
	CreateAttempt(ctx context.Context, q db.QueryExecutor, attempt *domain.QuizAttempt) (*domain.QuizAttempt, error)
	GetAttemptByID(ctx context.Context, q db.QueryExecutor, id int64) (*domain.QuizAttempt, error)
	UpdateAttempt(ctx context.Context, q db.QueryExecutor, attempt *domain.QuizAttempt) error
	CountUserAttempts(ctx context.Context, q db.QueryExecutor, userID, quizID int64) (int, error)
	CountUserAttemptsForUpdate(ctx context.Context, q db.QueryExecutor, userID, quizID int64) (int, error)
	CreateBatchAnswers(ctx context.Context, q db.QueryExecutor, answers []domain.QuizAttemptAnswer) error
	UpdateAttemptAnswer(ctx context.Context, q db.QueryExecutor, answerID int64, points int, feedback *string, isCorrect bool) error
	CountUngradedAnswers(ctx context.Context, q db.QueryExecutor, attemptID int64) (int, error)
	GetAnswerPointsAndCorrectness(ctx context.Context, q db.QueryExecutor, answerID int64) (bool, int, error)
	UpdateLessonProgressAfterQuiz(ctx context.Context, q db.QueryExecutor, userID, lessonID int64, score int) error
	SumAttemptPoints(ctx context.Context, q db.QueryExecutor, attemptID int64) (int, error)
	GetQuizTotalPoints(ctx context.Context, q db.QueryExecutor, quizID int64) (int, error)
	GetAttemptForUpdate(ctx context.Context, q db.QueryExecutor, attemptID int64) (*domain.QuizAttempt, error)
	AcquireAdvisoryLock(ctx context.Context, q db.QueryExecutor, userID int64, quizID int64) error
	ListAttemptsForGrading(ctx context.Context, q db.QueryExecutor, courseID int64, quizID *int64, limit, offset int) ([]domain.QuizAttempt, int64, error)
	SaveEssaySubmission(ctx context.Context, q db.QueryExecutor, userID, courseID, lessonID int64, essay domain.EssaySubmission) error
}
