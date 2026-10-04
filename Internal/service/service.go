package service

import (
	"context"
	"io"

	"edtech/internal/domain"
)

// --- Category Interfaces ---

type CategoryReader interface {
	GetCategoryByID(ctx context.Context, id int64) (*domain.Category, error)
	GetCategoryBySlug(ctx context.Context, slug string) (*domain.Category, error)
	ListCategories(ctx context.Context) ([]domain.Category, error)
}

type CategoryManager interface {
	CreateCategory(ctx context.Context, category *domain.Category) (int64, error)
	UpdateCategory(ctx context.Context, category *domain.Category) error
	DeleteCategory(ctx context.Context, id int64) error
}

type CategoryServices interface {
	CategoryReader
	CategoryManager
}

// --- Course Interfaces ---

// CourseReader defines read-only operations on course catalog and content structure
type CourseReader interface {
	GetCourseByID(ctx context.Context, id int64) (*domain.Course, error)
	GetCourseBySlug(ctx context.Context, slug string) (*domain.Course, error)
	GetCourseWithLessons(ctx context.Context, courseID int64) (domain.CourseWithLessons, error)
	GetCourseStructure(ctx context.Context, courseID int64) (domain.CourseStructure, error)
	ListPublicCourses(ctx context.Context, pagination domain.Pagination, filter domain.CourseFilter) (domain.PaginatedCourses, error)
}

// CourseManager defines authoring, state lifecycle and curriculum structuring
type CourseManager interface {
	CreateCourse(ctx context.Context, course *domain.Course) (*domain.Course, error)
	UpdateCourse(ctx context.Context, courseID int64, userID int64, input domain.UpdateCourseInput) (*domain.Course, error)
	UpdateCourseStatus(ctx context.Context, userID int64, courseID int64, status string) error
	PublishCourse(ctx context.Context, userID int64, courseID int64) error
	ArchiveCourse(ctx context.Context, userID int64, courseID int64) error
	DeleteCourse(ctx context.Context, userID int64, courseID int64) error
	ReorderSections(ctx context.Context, userID int64, courseID int64, sectionIDs []int64) error
}

// CourseStudent defines student course views
type CourseStudent interface {
	ListMyCourses(ctx context.Context, input *domain.InputListMyCourse) (domain.PaginatedCourses, error)
}

// CourseServices aggregates all course operations for backward compatibility and full service access
type CourseServices interface {
	CourseReader
	CourseManager
	CourseStudent
}

// --- Enrollment Interfaces ---

// StudentEnrollmentOperations handles self-enrollment and unenrollment by students
type StudentEnrollmentOperations interface {
	SelfEnrollCourse(ctx context.Context, req domain.SelfEnrollRequest) error
	UnenrollUser(ctx context.Context, userID int64, courseID int64) error
	GetRoleUserInCource(ctx context.Context, userID, courseID int64) (string, error)
}

// TeacherEnrollmentOperations handles course user roster management by instructors
type TeacherEnrollmentOperations interface {
	TeacherEnrollCourse(ctx context.Context, req domain.TeacherEnrollRequest) error
	TeacherUnenrollUser(ctx context.Context, teacherID int64, targetUserID int64, courseID int64) error
	ListCourseStudents(ctx context.Context, courseID int64, page int64, pageSize int64) ([]domain.User, int, error)
	ListCourseStudentsWithProgress(ctx context.Context, teacherID, courseID int64, page, pageSize int64) ([]domain.CourseStudentItem, int64, error)
	ChangeUserRole(ctx context.Context, courseID int64, targetUserID int64, newRole string) error
}

// EnrolledServices aggregates all enrollment operations
type EnrolledServices interface {
	StudentEnrollmentOperations
	TeacherEnrollmentOperations
}

// --- Progress Interfaces ---

// ProgressTrackingOperations handles student lesson progress and quiz completion recording
type ProgressTrackingOperations interface {
	StartLesson(ctx context.Context, userID int64, lessonID int64) error
	UpdateLessonProgress(ctx context.Context, userID int64, lessonID int64, input domain.UpdateProgressInput) error
	CompleteLesson(ctx context.Context, userID int64, lessonID int64, input domain.CompleteLessonInput) (*domain.LessonCompletionResult, error)
	StartLessonAttempt(ctx context.Context, userID int64, lessonID int64) (*domain.StartAttemptResult, error)
}

// ProgressReaderOperations handles querying lesson and course progress statistics
type ProgressReaderOperations interface {
	GetLessonProgress(ctx context.Context, userID int64, lessonID int64) (*domain.LessonProgress, error)
	GetCourseProgress(ctx context.Context, userID int64, courseID int64) (*domain.CourseProgress, error)
	GetAllLessonProgress(ctx context.Context, userID int64, courseID int64) ([]domain.LessonProgress, error)
	GetLessonAttemptsSummary(ctx context.Context, userID int64, lessonID int64) (*domain.LessonAttemptsSummary, error)
}

// ProgressServices aggregates all progress tracking operations
type ProgressServices interface {
	ProgressTrackingOperations
	ProgressReaderOperations
}

// --- Auth Interfaces ---

// AuthOperations handles authentication and token lifecycle
type AuthOperations interface {
	Register(ctx context.Context, input domain.RegisterInput) (*domain.User, error)
	Login(ctx context.Context, email, password string) (domain.AuthTokens, error)
	RefreshToken(ctx context.Context, refreshToken string) (domain.AuthTokens, error)
	Logout(ctx context.Context, refreshToken string) error
	VerifyEmail(ctx context.Context, userID int64) error
}

// UserProfileOperations handles user profile inspection and updates
type UserProfileOperations interface {
	GetCurrentUser(ctx context.Context, userID int64) (*domain.User, error)
	UpdateProfile(ctx context.Context, userID int64, input domain.UpdateProfileInput) (*domain.User, error)
	ChangePassword(ctx context.Context, userID int64, oldPassword, newPassword string) error
}

// UserAdminOperations handles administrative user management
type UserAdminOperations interface {
	ChangeUserRole(ctx context.Context, adminID int64, targetUserID int64, newRole domain.Role) error
	SetUserBanned(ctx context.Context, adminID int64, targetUserID int64, isBanned bool) error
	ListUsersForAdmin(ctx context.Context, adminID int64, filter domain.UserFilter, pagination domain.Pagination) ([]domain.User, int64, error)
}

// AuthService aggregates all user and authentication operations
type AuthService interface {
	AuthOperations
	UserProfileOperations
	UserAdminOperations
}

// --- Lesson Interfaces ---

type LessonServices interface {
	CreateLesson(ctx context.Context, userID int64, lesson *domain.Lesson) (int64, error)
	GetLesson(ctx context.Context, userID int64, lessonID int64) (*domain.Lesson, error)
	DeleteLesson(ctx context.Context, userID int64, lessonID int64) error
	UpdateLesson(ctx context.Context, userID int64, lesson *domain.Lesson) error
	UpdateLessonStatus(ctx context.Context, userID int64, lessonID int64, status string) error
}

// --- Analytics Interfaces ---

type AnalyticsServices interface {
	GetCourseAnalytics(ctx context.Context, teacherID, courseID int64) (domain.CourseAnalyticsSummary, error)
	GetStudentDrilldown(ctx context.Context, teacherID, courseID, studentID int64) (*domain.StudentDrilldownReport, error)
	ListPendingHomeworks(ctx context.Context, teacherID, courseID int64, page, pageSize int64) ([]domain.PendingHomeworkItem, int64, error)
	ListTeacherPendingHomeworks(ctx context.Context, teacherID int64, courseID int64, page, pageSize int64) (*domain.TeacherPendingHomeworksResult, error)
}

// --- Access Interfaces ---

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

// --- Quiz Interfaces ---

type QuizServices interface {
	CreateQuiz(ctx context.Context, userID int64, quiz *domain.Quiz) (*domain.Quiz, error)
	StartAttempt(ctx context.Context, userID int64, quizID int64) (*domain.QuizAttempt, error)
	SubmitAttempt(ctx context.Context, userID int64, attemptID int64, answers []domain.QuizAttemptAnswer) (*domain.QuizAttempt, error)
	ListAttemptsForGrading(ctx context.Context, userID int64, courseID int64, quizID *int64, page int, pageSize int) ([]domain.QuizAttempt, int64, error)
	GradeAttemptAnswer(ctx context.Context, teacherID int64, attemptID int64, answerID int64, points int, feedback *string) (*domain.QuizAttempt, error)
}

// --- Section Interfaces ---

type SectionServices interface {
	CreateSection(ctx context.Context, userID int64, section *domain.Section) (*domain.Section, error)
	UpdateSection(ctx context.Context, userID int64, section *domain.Section) error
	UpdateSectionStatus(ctx context.Context, userID int64, sectionID int64, status string) error
	ReorderLessons(ctx context.Context, userID int64, sectionID int64, lessonIDs []int64) error
	DeleteSection(ctx context.Context, userID int64, sectionID int64) error
}

// --- Upload Interfaces ---

type UploadServices interface {
	UploadFile(ctx context.Context, file io.Reader, filename string, size int64, category string) (*domain.FileUploadResult, error)
	UploadImagesBatch(ctx context.Context, files []domain.BatchFileItem, category string) ([]domain.BatchUploadResultItem, error)
}

// --- Review Interfaces ---

type ReviewServices interface {
	AddOrUpdateReview(ctx context.Context, userID, courseID int64, rating int, comment *string) (*domain.Review, error)
	DeleteReview(ctx context.Context, userID, courseID int64) error
	ListCourseReviews(ctx context.Context, courseID int64, page, pageSize int) (*domain.CourseReviewsSummary, error)
	GetMyReview(ctx context.Context, userID, courseID int64) (*domain.Review, error)
}

// --- Certificate Interfaces ---

type CertificateServices interface {
	GetOrIssueCertificate(ctx context.Context, userID, courseID int64) (*domain.Certificate, error)
	VerifyCertificate(ctx context.Context, code string) (*domain.Certificate, error)
}

// --- Notification Interfaces ---

type NotificationServices interface {
	CreateNotification(ctx context.Context, userID int64, title, message string, nType domain.NotificationType, linkURL *string) (*domain.Notification, error)
	GetFeed(ctx context.Context, userID int64, limit int) (*domain.NotificationFeed, error)
	MarkAsRead(ctx context.Context, userID, notificationID int64) error
	MarkAllAsRead(ctx context.Context, userID int64) error
}
