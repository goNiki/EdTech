package repository

import (
	"context"
	"edtech/internal/domain"
	"time"
)

type UserRepository interface {
	CreateUser(ctx context.Context, createUser domain.CreateUser) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserByID(ctx context.Context, userID int64) (*domain.User, error)
	GetUserByUserName(ctx context.Context, username string) (*domain.User, error)
	ExistingByEmail(ctx context.Context, email string) (bool, error)
	ExistingByUsernName(ctx context.Context, username string) (bool, error)
	UpdateLastLogin(ctx context.Context, userID int64, now time.Time) error
	UpdateProfile(ctx context.Context, user *domain.User) error
	UpdatePassword(ctx context.Context, userID int64, passHash string) error
	SetEmailVerified(ctx context.Context, userID int64, verified bool) error
	UpdateRole(ctx context.Context, userID int64, role domain.Role) error
	SetBannedStatus(ctx context.Context, userID int64, isBanned bool) error
	ListUsers(ctx context.Context, filter domain.UserFilter, pagination domain.Pagination) ([]domain.User, int64, error)
}

type CategoryRepository interface {
	ListCategories(ctx context.Context) ([]domain.Category, error)
	CreateCategory(ctx context.Context, category *domain.Category) (*domain.Category, error)
	GetCategoryByID(ctx context.Context, id int64) (*domain.Category, error)
	GetCategoryBySlug(ctx context.Context, slug string) (*domain.Category, error)
	UpdateCategory(ctx context.Context, category *domain.Category) error
	DeleteCategory(ctx context.Context, id int64) error
}

type CourseRepository interface {
	GetCourseBySlug(ctx context.Context, slug string) (*domain.Course, error)
	CreateCourse(ctx context.Context, course *domain.Course) (*domain.Course, error)
	GetCourseByID(ctx context.Context, id int64) (*domain.Course, error)
	ArchiveCourse(ctx context.Context, course *domain.Course) error
	PublishCourse(ctx context.Context, course *domain.Course) error
	UpdateCourseStatus(ctx context.Context, courseID int64, status string) error
	UpdateCourse(ctx context.Context, course *domain.Course) error
	ListPublicCourses(ctx context.Context, pagination domain.Pagination, filter domain.CourseFilter) ([]domain.Course, error)
	ListEnrolledCourses(ctx context.Context, input *domain.InputListMyCourse) ([]domain.Course, error)
	CountEnrolledCourses(ctx context.Context, input *domain.InputListMyCourse) (int64, error)
	CountCourses(ctx context.Context, filter domain.CourseFilter) (int64, error)
	DeleteCourse(ctx context.Context, courseID int64) error
	ExistingBySlug(ctx context.Context, slug string) (bool, error)
	ExistsBySlug(ctx context.Context, slug string) (bool, error)
	IncrementEnrolledCount(ctx context.Context, courseID int64) error
	DecrementEnrolledCount(ctx context.Context, courseID int64) error
	UpdateCourseRatingStats(ctx context.Context, courseID int64, rating float64, reviewsCount int) error
}

type EnrolledRepository interface {
	UserExistCourse(ctx context.Context, userID, courseid int64) (bool, error)
	EnrollUserToCourse(ctx context.Context, enroll domain.EnrolledInCourse) error
	GetRoleUserInCourse(ctx context.Context, userID, courceID int64) (string, error)
	UnenrollUser(ctx context.Context, userID, courseID int64) error
	ListCourseStudents(ctx context.Context, courseID int64, limit, offset int64) ([]domain.User, int, error)
	ListCourseStudentsWithProgress(ctx context.Context, courseID int64, limit, offset int64) ([]domain.CourseStudentItem, int64, error)
	ChangeUserRole(ctx context.Context, courseID int64, targetUserID int64, newRole string) error
}

type SectionRepository interface {
	CreateSection(ctx context.Context, section *domain.Section) (*domain.Section, error)
	GetSectionByID(ctx context.Context, id int64) (*domain.Section, error)
	ListSectionsByCourseID(ctx context.Context, courseID int64) ([]domain.Section, error)
	UpdateSection(ctx context.Context, section *domain.Section) error
	UpdateSectionStatus(ctx context.Context, sectionID int64, status string) error
	UpdateStatusByCourseID(ctx context.Context, courseID int64, status string) error
	ReorderSections(ctx context.Context, courseID int64, sectionIDs []int64) error
	DeleteSection(ctx context.Context, sectionID int64) error
	GetMaxPositionByCourseID(ctx context.Context, courseID int64) (int, error)
}

type LessonRepository interface {
	GetMaxPositionByCourseID(ctx context.Context, courseId int64) (int64, error)
	CreateLesson(ctx context.Context, lesson *domain.Lesson) error
	GetLessonByID(ctx context.Context, id int64) (*domain.Lesson, error)
	UpdateLesson(ctx context.Context, lesson *domain.Lesson) error
	UpdateLessonStatus(ctx context.Context, lessonID int64, status string) error
	UpdateStatusByCourseID(ctx context.Context, courseID int64, status string) error
	UpdateStatusBySectionID(ctx context.Context, sectionID int64, status string) error
	ReorderLessons(ctx context.Context, sectionID *int64, lessonIDs []int64) error
	DeleteLessonByID(ctx context.Context, lessonID int64) error
	GetLessonsByCourseID(ctx context.Context, courseID int64) ([]domain.Lesson, error)
}

type AnalyticsRepository interface {
	GetCourseAnalyticsSummary(ctx context.Context, courseID int64) (domain.CourseAnalyticsSummary, error)
	ListPendingHomeworks(ctx context.Context, courseID int64, limit, offset int64) ([]domain.PendingHomeworkItem, int64, error)
	ListTeacherPendingHomeworks(ctx context.Context, teacherID int64, courseID int64, limit, offset int64) ([]domain.PendingHomeworkItem, int64, []domain.CoursePendingSummaryItem, error)
	GetStudentDrilldown(ctx context.Context, userID, courseID int64) (*domain.StudentDrilldownReport, error)
}

type RefreshRepository interface {
	Save(ctx context.Context, reftoken string, id int64, expiresAt time.Time) error
	GetByToken(ctx context.Context, refreshToken string) (domain.RefreshTokenData, error)
	DeleteRefreshToken(ctx context.Context, hashToken string) error
	DeleteAllByUserID(ctx context.Context, userID int64) error
}

type PermissionsRepository interface {
	HasPermission(ctx context.Context, roleName string, resource string, action string) (bool, error)
	GetRolePermissions(ctx context.Context, rolename string) ([]string, error)
}

type ProgressRepository interface {
	CreateLessonProgress(ctx context.Context, progress *domain.LessonProgress) error
	GetLessonProgress(ctx context.Context, userID, lessonID int64) (*domain.LessonProgress, error)
	UpdateLessonProgressTime(ctx context.Context, userID, lessonID int64, additionalTime int, lastPos int) error
	UpdateLessonProgressStatus(ctx context.Context, userID, lessonID int64, status domain.ProgressStatus) error
	UpdateLessonProgressScore(ctx context.Context, userID, lessonID int64, score int) error

	CreateCourseProgress(ctx context.Context, progress *domain.CourseProgress) error
	UpsertCourseProgress(ctx context.Context, userID, courseID int64, completedLessons, totalLessons int, percentage float64) error
	UpsertCourseProgressWithScore(ctx context.Context, userID, courseID int64, completedLessons, totalLessons int, percentage float64, averageScore float64) error
	GetCourseProgress(ctx context.Context, userID, courseID int64) (*domain.CourseProgress, error)
	GetAllLessonProgressByCourse(ctx context.Context, userID, courseID int64) ([]domain.LessonProgress, error)
}

type QuizRepository interface {
	CreateQuiz(ctx context.Context, quiz *domain.Quiz) (*domain.Quiz, error)
	GetQuizByID(ctx context.Context, id int64) (*domain.Quiz, error)
	GetQuizByLessonID(ctx context.Context, lessonID int64) (*domain.Quiz, error)
	CreateAttempt(ctx context.Context, attempt *domain.QuizAttempt) (*domain.QuizAttempt, error)
	GetAttemptByID(ctx context.Context, id int64) (*domain.QuizAttempt, error)
	UpdateAttempt(ctx context.Context, attempt *domain.QuizAttempt) error
	CountUserAttempts(ctx context.Context, userID, quizID int64) (int, error)
	CountUserAttemptsForUpdate(ctx context.Context, userID, quizID int64) (int, error)
	CreateBatchAnswers(ctx context.Context, answers []domain.QuizAttemptAnswer) error
	UpdateAttemptAnswer(ctx context.Context, answerID int64, points int, feedback *string, isCorrect bool, gradedBy ...int64) error
	CountUngradedAnswers(ctx context.Context, attemptID int64) (int, error)
	GetAnswerPointsAndCorrectness(ctx context.Context, answerID int64) (bool, int, error)
	UpdateLessonProgressAfterQuiz(ctx context.Context, userID, lessonID int64, score int) error
	SumAttemptPoints(ctx context.Context, attemptID int64) (int, error)
	GetQuizTotalPoints(ctx context.Context, quizID int64) (int, error)
	GetAttemptForUpdate(ctx context.Context, attemptID int64) (*domain.QuizAttempt, error)
	AcquireAdvisoryLock(ctx context.Context, userID int64, quizID int64) error
	ListAttemptsForGrading(ctx context.Context, courseID int64, quizID *int64, limit, offset int) ([]domain.QuizAttempt, int64, error)
	SaveEssaySubmission(ctx context.Context, userID, courseID, lessonID int64, essay domain.EssaySubmission) error
	GetLessonSubmissions(ctx context.Context, userID, lessonID int64) ([]domain.LessonSubmissionDetail, error)
	ListUserAttemptsByLessonID(ctx context.Context, userID, lessonID int64) ([]domain.QuizAttempt, error)
	GetBestScoreByLessonID(ctx context.Context, userID, lessonID int64) (int, error)
	SaveAttemptDraft(ctx context.Context, attemptID, userID int64, currentStep int, draftAnswers []byte) error
	GetActiveAttempt(ctx context.Context, userID, lessonID int64) (*domain.QuizAttempt, error)
	GetStudentHomeworkFeedback(ctx context.Context, userID, lessonID int64) (*domain.StudentHomeworkFeedback, error)
}


type ReviewRepository interface {
	UpsertReview(ctx context.Context, review *domain.Review) (*domain.Review, error)
	DeleteReview(ctx context.Context, courseID, userID int64) error
	GetReviewByUserAndCourse(ctx context.Context, courseID, userID int64) (*domain.Review, error)
	ListReviewsByCourse(ctx context.Context, courseID int64, limit, offset int) ([]domain.Review, int64, error)
	GetCourseRatingSummary(ctx context.Context, courseID int64) (avgRating float64, count int, err error)
}

type CertificateRepository interface {
	CreateCertificate(ctx context.Context, cert *domain.Certificate) (*domain.Certificate, error)
	GetCertificateByCode(ctx context.Context, code string) (*domain.Certificate, error)
	GetCertificateByUserAndCourse(ctx context.Context, userID, courseID int64) (*domain.Certificate, error)
}

type NotificationRepository interface {
	CreateNotification(ctx context.Context, n *domain.Notification) (*domain.Notification, error)
	GetUserNotifications(ctx context.Context, userID int64, limit int) ([]domain.Notification, error)
	CountUnreadNotifications(ctx context.Context, userID int64) (int64, error)
	MarkAsRead(ctx context.Context, userID, notificationID int64) error
	MarkAllAsRead(ctx context.Context, userID int64) error
}
