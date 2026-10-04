package app

import (
	"log/slog"
	"net/http"
	"os"

	"edtech/internal/infrastructure/config"
	"edtech/internal/infrastructure/db"
	"edtech/internal/infrastructure/hasher"
	"edtech/internal/infrastructure/jwt"
	"edtech/internal/infrastructure/logger"
	"edtech/internal/infrastructure/txmanager"

	analyticsHandler "edtech/internal/interfaces/handlers/analytics"
	authHandler "edtech/internal/interfaces/handlers/auth"
	courseHandler "edtech/internal/interfaces/handlers/courses"
	enrollmentHandler "edtech/internal/interfaces/handlers/enrollment"
	lessonHandler "edtech/internal/interfaces/handlers/lesson"
	progressHandler "edtech/internal/interfaces/handlers/progress"
	quizHandler "edtech/internal/interfaces/handlers/quiz"
	sectionHandler "edtech/internal/interfaces/handlers/section"
	uploadHandler "edtech/internal/interfaces/handlers/upload"
	categoryHandler "edtech/internal/interfaces/handlers/category"
	reviewHandler "edtech/internal/interfaces/handlers/review"
	certHandler "edtech/internal/interfaces/handlers/certificate"
	notifHandler "edtech/internal/interfaces/handlers/notification"
	mwauth "edtech/internal/interfaces/middleware/auth"
	mwlogger "edtech/internal/interfaces/middleware/logger"
	"edtech/internal/infrastructure/storage"

	"edtech/internal/repository"
	analyticsRepo "edtech/internal/repository/analytics"
	authRepo "edtech/internal/repository/auth"
	categoryRepo "edtech/internal/repository/category"
	certRepo "edtech/internal/repository/certificate"
	notifRepo "edtech/internal/repository/notification"
	courseRepo "edtech/internal/repository/course"
	enrolledRepo "edtech/internal/repository/enrollment"
	lessonRepo "edtech/internal/repository/lesson"
	permissionRepo "edtech/internal/repository/permission"
	progressRepo "edtech/internal/repository/progress"
	quizRepo "edtech/internal/repository/quiz"
	refreshRepo "edtech/internal/repository/refresh"
	reviewRepo "edtech/internal/repository/review"
	sectionRepo "edtech/internal/repository/section"

	"edtech/internal/service"
	accessService "edtech/internal/service/access"
	analyticsService "edtech/internal/service/analytics"
	authService "edtech/internal/service/auth"
	categoryService "edtech/internal/service/category"
	certService "edtech/internal/service/certificate"
	notifService "edtech/internal/service/notification"
	courseService "edtech/internal/service/course"
	enrolledService "edtech/internal/service/enrollment"
	lessonService "edtech/internal/service/lesson"
	progressService "edtech/internal/service/progress"
	quizService "edtech/internal/service/quiz"
	reviewService "edtech/internal/service/review"
	sectionService "edtech/internal/service/section"
	uploadService "edtech/internal/service/upload"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-playground/validator/v10"
)

const configPath = ".env"

type diContainer struct {
	// infrastructure
	loggerCfg   config.Logger
	postgresCfg config.Postgres
	serverCfg   config.Server
	jwtCfg      config.JWT
	logger      *slog.Logger
	db          *db.Postgres
	txManager   *txmanager.TxManager
	hasher      hasher.HasherManager
	jwtManager  *jwt.JwtManager
	storage     storage.Storage

	// middleware
	mwAuth mwauth.AuthMiddleware
	mwLog  func(http.Handler) http.Handler

	// router
	router http.Handler

	// handlers
	authHdl       *authHandler.AuthHandler
	courseHdl     *courseHandler.CourseHandler
	lessonHdl     *lessonHandler.LessonHandler
	sectionHdl    *sectionHandler.SectionHandler
	uploadHdl     *uploadHandler.UploadHandler
	progressHdl   *progressHandler.ProgressHandler
	quizHdl       *quizHandler.QuizHandler
	enrollmentHdl *enrollmentHandler.EnrollmentHandler
	analyticsHdl  *analyticsHandler.AnalyticsHandler
	categoryHdl   *categoryHandler.CategoryHandler
	reviewHdl     *reviewHandler.ReviewHandler
	certHdl       *certHandler.CertificateHandler
	notifHdl      *notifHandler.NotificationHandler

	// services
	accessSvc    service.AccessService
	authSvc      service.AuthService
	courseSvc    service.CourseServices
	enrolledSvc  service.EnrolledServices
	lessonSvc    service.LessonServices
	sectionSvc   service.SectionServices
	uploadSvc    service.UploadServices
	progressSvc  service.ProgressServices
	quizSvc      service.QuizServices
	analyticsSvc service.AnalyticsServices
	categorySvc  service.CategoryServices
	reviewSvc    service.ReviewServices
	certSvc      service.CertificateServices
	notifSvc     service.NotificationServices

	// repositories
	userRepo      repository.UserRepository
	refreshRepo   repository.RefreshRepository
	courseRepo    repository.CourseRepository
	sectionRepo   repository.SectionRepository
	lessonRepo    repository.LessonRepository
	enrolledRepo  repository.EnrolledRepository
	permRepo      repository.PermissionsRepository
	progRepo      repository.ProgressRepository
	quizRepo      repository.QuizRepository
	analyticsRepo repository.AnalyticsRepository
	categoryRepo  repository.CategoryRepository
	reviewRepo    repository.ReviewRepository
	certRepo      repository.CertificateRepository
	notifRepo     repository.NotificationRepository
}

func (d *diContainer) initConfig() {
	if d.loggerCfg != nil && d.postgresCfg != nil && d.serverCfg != nil && d.jwtCfg != nil {
		return
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		slog.Error("failed to get config: ", "error", err)
		os.Exit(1)
	}

	d.loggerCfg = cfg.Logger
	d.postgresCfg = cfg.Postgres
	d.serverCfg = cfg.Server
	d.jwtCfg = cfg.JWT
}

func (d *diContainer) LoggerCfg() config.Logger {
	d.initConfig()
	return d.loggerCfg
}

func (d *diContainer) PostgresCfg() config.Postgres {
	d.initConfig()
	return d.postgresCfg
}

func (d *diContainer) ServerCfg() config.Server {
	d.initConfig()
	return d.serverCfg
}

func (d *diContainer) JwtCfg() config.JWT {
	d.initConfig()
	return d.jwtCfg
}

func (d *diContainer) Logger() *slog.Logger {
	if d.logger == nil {
		d.logger = logger.Init()
	}
	return d.logger
}

func (d *diContainer) DB() *db.Postgres {
	if d.db == nil {
		var err error
		d.db, err = db.New(d.PostgresCfg())
		if err != nil {
			slog.Error("failed to create database connection: ", "error", err)
			os.Exit(1)
		}
	}
	return d.db
}

func (d *diContainer) TxManager() *txmanager.TxManager {
	if d.txManager == nil {
		d.txManager = txmanager.NewTxManager(d.DB())
	}
	return d.txManager
}

func (d *diContainer) Hasher() hasher.HasherManager {
	if d.hasher == nil {
		d.hasher = hasher.NewHasher()
	}
	return d.hasher
}

func (d *diContainer) JWTManager() *jwt.JwtManager {
	if d.jwtManager == nil {
		d.jwtManager = jwt.NewJwtManager(d.JwtCfg())
	}
	return d.jwtManager
}

func (d *diContainer) Storage() storage.Storage {
	if d.storage == nil {
		var err error
		d.storage, err = storage.NewLocalStorage("./uploads")
		if err != nil {
			slog.Error("failed to create storage: ", "error", err)
			os.Exit(1)
		}
	}
	return d.storage
}

func (d *diContainer) Close() {
	if d.db != nil && d.db.Pool != nil {
		d.db.Pool.Close()
	}
}

// Middleware

func (d *diContainer) MwAuth() mwauth.AuthMiddleware {
	if d.mwAuth == nil {
		d.mwAuth = mwauth.NewAuthMiddleware(d.JWTManager())
	}
	return d.mwAuth
}

func (d *diContainer) MwLog() func(http.Handler) http.Handler {
	if d.mwLog == nil {
		d.mwLog = mwlogger.NewLoggerMiddleware(d.Logger()).Log
	}
	return d.mwLog
}

// Repositories

func (d *diContainer) UserRepo() repository.UserRepository {
	if d.userRepo == nil {
		d.userRepo = authRepo.NewAuthRepo(d.DB().Pool)
	}
	return d.userRepo
}

func (d *diContainer) RefreshRepo() repository.RefreshRepository {
	if d.refreshRepo == nil {
		d.refreshRepo = refreshRepo.NewRefreshRepo(d.DB().Pool)
	}
	return d.refreshRepo
}

func (d *diContainer) CourseRepo() repository.CourseRepository {
	if d.courseRepo == nil {
		d.courseRepo = courseRepo.NewCourseRepo(d.DB().Pool)
	}
	return d.courseRepo
}

func (d *diContainer) SectionRepo() repository.SectionRepository {
	if d.sectionRepo == nil {
		d.sectionRepo = sectionRepo.NewSectionRepository(d.DB().Pool)
	}
	return d.sectionRepo
}

func (d *diContainer) LessonRepo() repository.LessonRepository {
	if d.lessonRepo == nil {
		d.lessonRepo = lessonRepo.NewLessonRepo(d.DB().Pool)
	}
	return d.lessonRepo
}

func (d *diContainer) EnrolledRepo() repository.EnrolledRepository {
	if d.enrolledRepo == nil {
		d.enrolledRepo = enrolledRepo.NewEnrolledRepo(d.DB().Pool)
	}
	return d.enrolledRepo
}

func (d *diContainer) PermRepo() repository.PermissionsRepository {
	if d.permRepo == nil {
		d.permRepo = permissionRepo.NewPermissionRepo(d.DB().Pool)
	}
	return d.permRepo
}

func (d *diContainer) ProgRepo() repository.ProgressRepository {
	if d.progRepo == nil {
		d.progRepo = progressRepo.NewProgressRepo(d.DB().Pool)
	}
	return d.progRepo
}

func (d *diContainer) QuizRepo() repository.QuizRepository {
	if d.quizRepo == nil {
		d.quizRepo = quizRepo.NewQuizRepo(d.DB().Pool)
	}
	return d.quizRepo
}

func (d *diContainer) AnalyticsRepo() repository.AnalyticsRepository {
	if d.analyticsRepo == nil {
		d.analyticsRepo = analyticsRepo.NewAnalyticsRepo(d.DB().Pool)
	}
	return d.analyticsRepo
}

func (d *diContainer) CategoryRepo() repository.CategoryRepository {
	if d.categoryRepo == nil {
		d.categoryRepo = categoryRepo.NewCategoryRepo(d.DB().Pool)
	}
	return d.categoryRepo
}

func (d *diContainer) ReviewRepo() repository.ReviewRepository {
	if d.reviewRepo == nil {
		d.reviewRepo = reviewRepo.NewReviewRepository(d.DB().Pool)
	}
	return d.reviewRepo
}

func (d *diContainer) CertRepo() repository.CertificateRepository {
	if d.certRepo == nil {
		d.certRepo = certRepo.NewCertificateRepository(d.DB().Pool)
	}
	return d.certRepo
}

func (d *diContainer) NotifRepo() repository.NotificationRepository {
	if d.notifRepo == nil {
		d.notifRepo = notifRepo.NewNotificationRepository(d.DB().Pool)
	}
	return d.notifRepo
}

// Services

func (d *diContainer) SectionSvc() service.SectionServices {
	if d.sectionSvc == nil {
		d.sectionSvc = sectionService.NewSectionService(d.SectionRepo(), d.LessonRepo(), d.CourseRepo(), d.AccessSvc(), d.TxManager())
	}
	return d.sectionSvc
}

func (d *diContainer) AccessSvc() service.AccessService {
	if d.accessSvc == nil {
		d.accessSvc = accessService.NewAccessService(d.CourseRepo(), d.EnrolledRepo(), d.PermRepo())
	}
	return d.accessSvc
}

func (d *diContainer) AuthSvc() service.AuthService {
	if d.authSvc == nil {
		d.authSvc = authService.NewAuthService(d.UserRepo(), d.JWTManager(), d.Hasher(), d.RefreshRepo(), d.TxManager())
	}
	return d.authSvc
}

func (d *diContainer) CourseSvc() service.CourseServices {
	if d.courseSvc == nil {
		d.courseSvc = courseService.NewCourseService(d.CourseRepo(), d.SectionRepo(), d.LessonRepo(), d.AccessSvc(), d.EnrolledRepo(), d.TxManager())
	}
	return d.courseSvc
}

func (d *diContainer) LessonSvc() service.LessonServices {
	if d.lessonSvc == nil {
		d.lessonSvc = lessonService.NewLessonService(d.CourseRepo(), d.LessonRepo(), d.SectionRepo(), d.AccessSvc())
	}
	return d.lessonSvc
}

func (d *diContainer) EnrolledSvc() service.EnrolledServices {
	if d.enrolledSvc == nil {
		d.enrolledSvc = enrolledService.NewEnrolmentService(d.CourseRepo(), d.EnrolledRepo(), d.UserRepo(), d.AccessSvc(), d.TxManager())
	}
	return d.enrolledSvc
}

func (d *diContainer) ProgressSvc() service.ProgressServices {
	if d.progressSvc == nil {
		d.progressSvc = progressService.NewProgressService(d.ProgRepo(), d.LessonRepo(), d.QuizRepo(), d.TxManager())
	}
	return d.progressSvc
}

func (d *diContainer) QuizSvc() service.QuizServices {
	if d.quizSvc == nil {
		d.quizSvc = quizService.NewQuizService(d.QuizRepo(), d.CourseRepo(), d.LessonRepo(), d.AccessSvc(), d.ProgressSvc(), d.TxManager(), d.NotifSvc())
	}
	return d.quizSvc
}

func (d *diContainer) AnalyticsSvc() service.AnalyticsServices {
	if d.analyticsSvc == nil {
		d.analyticsSvc = analyticsService.NewAnalyticsService(d.AnalyticsRepo(), d.CourseRepo(), d.AccessSvc(), d.TxManager())
	}
	return d.analyticsSvc
}

func (d *diContainer) UploadSvc() service.UploadServices {
	if d.uploadSvc == nil {
		d.uploadSvc = uploadService.NewUploadService(d.Storage())
	}
	return d.uploadSvc
}

func (d *diContainer) CategorySvc() service.CategoryServices {
	if d.categorySvc == nil {
		d.categorySvc = categoryService.NewCategoryService(d.CategoryRepo())
	}
	return d.categorySvc
}

func (d *diContainer) ReviewSvc() service.ReviewServices {
	if d.reviewSvc == nil {
		d.reviewSvc = reviewService.NewReviewService(d.ReviewRepo(), d.CourseRepo(), d.EnrolledRepo(), d.ProgRepo(), d.TxManager())
	}
	return d.reviewSvc
}

func (d *diContainer) CertSvc() service.CertificateServices {
	if d.certSvc == nil {
		d.certSvc = certService.NewCertificateService(d.CertRepo(), d.CourseRepo(), d.UserRepo(), d.ProgRepo())
	}
	return d.certSvc
}

func (d *diContainer) NotifSvc() service.NotificationServices {
	if d.notifSvc == nil {
		d.notifSvc = notifService.NewNotificationService(d.NotifRepo())
	}
	return d.notifSvc
}

// Handlers

func (d *diContainer) LessonHdl() *lessonHandler.LessonHandler {
	if d.lessonHdl == nil {
		d.lessonHdl = lessonHandler.NewLessonHandler(d.LessonSvc(), d.Logger(), validator.New(), d.MwAuth())
	}
	return d.lessonHdl
}

func (d *diContainer) SectionHdl() *sectionHandler.SectionHandler {
	if d.sectionHdl == nil {
		d.sectionHdl = sectionHandler.NewSectionHandler(d.SectionSvc(), d.Logger(), validator.New(), d.MwAuth())
	}
	return d.sectionHdl
}

func (d *diContainer) UploadHdl() *uploadHandler.UploadHandler {
	if d.uploadHdl == nil {
		d.uploadHdl = uploadHandler.NewUploadHandler(d.UploadSvc(), d.MwAuth(), d.Logger())
	}
	return d.uploadHdl
}

func (d *diContainer) AuthHdl() *authHandler.AuthHandler {
	if d.authHdl == nil {
		d.authHdl = authHandler.NewAuthHandler(d.AuthSvc(), d.MwAuth())
	}
	return d.authHdl
}

func (d *diContainer) CourseHdl() *courseHandler.CourseHandler {
	if d.courseHdl == nil {
		d.courseHdl = courseHandler.NewCourseHandler(d.CourseSvc(), d.LessonSvc(), d.EnrolledSvc(), d.MwAuth(), d.AccessSvc())
	}
	return d.courseHdl
}

func (d *diContainer) ProgressHdl() *progressHandler.ProgressHandler {
	if d.progressHdl == nil {
		d.progressHdl = progressHandler.NewProgressHandler(d.ProgressSvc(), d.MwAuth())
	}
	return d.progressHdl
}

func (d *diContainer) QuizHdl() *quizHandler.QuizHandler {
	if d.quizHdl == nil {
		d.quizHdl = quizHandler.NewQuizHandler(d.QuizSvc(), d.MwAuth())
	}
	return d.quizHdl
}

func (d *diContainer) EnrollmentHdl() *enrollmentHandler.EnrollmentHandler {
	if d.enrollmentHdl == nil {
		d.enrollmentHdl = enrollmentHandler.NewEnrollmentHandler(d.EnrolledSvc(), d.MwAuth())
	}
	return d.enrollmentHdl
}

func (d *diContainer) AnalyticsHdl() *analyticsHandler.AnalyticsHandler {
	if d.analyticsHdl == nil {
		d.analyticsHdl = analyticsHandler.NewAnalyticsHandler(d.AnalyticsSvc(), d.MwAuth())
	}
	return d.analyticsHdl
}

func (d *diContainer) CategoryHdl() *categoryHandler.CategoryHandler {
	if d.categoryHdl == nil {
		d.categoryHdl = categoryHandler.NewCategoryHandler(d.CategorySvc(), d.MwAuth())
	}
	return d.categoryHdl
}

func (d *diContainer) ReviewHdl() *reviewHandler.ReviewHandler {
	if d.reviewHdl == nil {
		d.reviewHdl = reviewHandler.NewReviewHandler(d.ReviewSvc(), d.Logger(), validator.New(), d.MwAuth())
	}
	return d.reviewHdl
}

func (d *diContainer) CertHdl() *certHandler.CertificateHandler {
	if d.certHdl == nil {
		d.certHdl = certHandler.NewCertificateHandler(d.CertSvc(), d.Logger(), d.MwAuth())
	}
	return d.certHdl
}

func (d *diContainer) NotifHdl() *notifHandler.NotificationHandler {
	if d.notifHdl == nil {
		d.notifHdl = notifHandler.NewNotificationHandler(d.NotifSvc(), d.Logger(), d.MwAuth())
	}
	return d.notifHdl
}

// Router

func (d *diContainer) Router() http.Handler {
	if d.router == nil {
		r := chi.NewRouter()

		// CORS middleware
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins:   d.ServerCfg().CorsAllowedOrigins(),
			AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300,
		}))

		r.Use(middleware.RequestID)
		r.Use(d.MwLog())
		r.Use(middleware.Recoverer)

		// Auth
		r.Route("/api/v1/auth", func(r chi.Router) {
			r.Post("/register", d.AuthHdl().Register)
			r.Post("/login", d.AuthHdl().Login)
			r.Post("/refresh", d.AuthHdl().RefreshToken)
			r.Post("/logout", d.AuthHdl().Logout)

			r.Group(func(r chi.Router) {
				r.Use(d.MwAuth().JWTMiddleware)

				r.Get("/me", d.AuthHdl().GetCurrentUser)
				r.Patch("/profile", d.AuthHdl().UpdateProfile)
				r.Post("/change-password", d.AuthHdl().ChangePassword)
				r.Post("/verify-email", d.AuthHdl().VerifyEmail)
			})
		})

		// Courses
		r.Route("/api/v1/courses", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(d.MwAuth().OptionalJWTMiddleware)
				r.Get("/", d.CourseHdl().ListPublicCourses)
				r.Get("/{courseid}", d.CourseHdl().GetCourseByID)
				r.Get("/slug/{slug}", d.CourseHdl().GetCourseBySlug)
				r.Get("/{courseid}/structure", d.CourseHdl().GetCourseStructure)
				r.Get("/{courseid}/reviews", d.ReviewHdl().ListReviews)
			})

			r.Group(func(r chi.Router) {
				r.Use(d.MwAuth().JWTMiddleware)

				r.Get("/my", d.CourseHdl().ListMyCourses)
				r.Post("/", d.CourseHdl().CreateCourse)
				r.Patch("/{courseid}", d.CourseHdl().UpdateCourse)
				r.Delete("/{courseid}", d.CourseHdl().DeleteCourse)
				r.Post("/{courseid}/publish", d.CourseHdl().PublishCourse)
				r.Post("/{courseid}/archive", d.CourseHdl().ArchiveCourse)
				r.Patch("/{courseid}/status", d.CourseHdl().UpdateCourseStatus)
				r.Put("/{courseid}/reorder-sections", d.CourseHdl().ReorderSections)

				// Reviews
				r.Post("/{courseid}/reviews", d.ReviewHdl().AddOrUpdateReview)
				r.Delete("/{courseid}/reviews", d.ReviewHdl().DeleteReview)
				r.Get("/{courseid}/reviews/my", d.ReviewHdl().GetMyReview)

				// Certificates
				r.Get("/{courseid}/certificate", d.CertHdl().GetOrIssueCertificate)

				// Enrollment
				r.Post("/{courseid}/enroll", d.EnrollmentHdl().SelfEnroll)
				r.Delete("/{courseid}/enroll", d.EnrollmentHdl().UnenrollSelf)
				r.Get("/{courseid}/students", d.EnrollmentHdl().ListCourseStudents)
				r.Post("/{courseid}/students", d.EnrollmentHdl().TeacherEnroll)
				r.Delete("/{courseid}/students/{userid}", d.EnrollmentHdl().TeacherUnenroll)

				// Analytics & Grading
				r.Get("/{courseid}/analytics", d.AnalyticsHdl().GetCourseAnalytics)
				r.Get("/{courseid}/students/{userid}/drilldown", d.AnalyticsHdl().GetStudentDrilldown)
				r.Get("/{courseid}/grading/pending", d.AnalyticsHdl().ListPendingHomeworks)

				// Progress & Quizzes
				r.Get("/{course_id}/progress", d.ProgressHdl().GetCourseProgress)
				r.Get("/{course_id}/progress/lessons", d.ProgressHdl().GetAllLessonProgress)
				r.Get("/{course_id}/quizzes/attempts", d.QuizHdl().ListAttemptsForGrading)
			})
		})

		// Sections
		r.Route("/api/v1/sections", func(r chi.Router) {
			r.Use(d.MwAuth().JWTMiddleware)

			r.Post("/", d.SectionHdl().CreateSection)
			r.Patch("/{id}", d.SectionHdl().UpdateSection)
			r.Patch("/{id}/status", d.SectionHdl().UpdateSectionStatus)
			r.Put("/{id}/reorder-lessons", d.SectionHdl().ReorderLessons)
			r.Delete("/{id}", d.SectionHdl().DeleteSection)
		})

		// Lessons
		r.Route("/api/v1/lessons", func(r chi.Router) {
			r.Use(d.MwAuth().JWTMiddleware)

			r.Post("/", d.LessonHdl().CreateLesson)
			r.Get("/{id}", d.LessonHdl().GetLesson)
			r.Patch("/{id}", d.LessonHdl().UpdateLesson)
			r.Patch("/{id}/status", d.LessonHdl().UpdateLessonStatus)
			r.Delete("/{id}", d.LessonHdl().DeleteLesson)

			r.Post("/{lesson_id}/start", d.ProgressHdl().StartLesson)
			r.Patch("/{lesson_id}/progress", d.ProgressHdl().UpdateLessonProgress)
			r.Post("/{lesson_id}/complete", d.ProgressHdl().CompleteLesson)
			r.Get("/{lesson_id}/progress", d.ProgressHdl().GetLessonProgress)
			r.Post("/{lesson_id}/quizzes", d.QuizHdl().CreateQuiz)
		})

		// Quizzes
		r.Route("/api/v1/quizzes", func(r chi.Router) {
			r.Use(d.MwAuth().JWTMiddleware)

			r.Post("/{quiz_id}/attempts/start", d.QuizHdl().StartAttempt)
			r.Post("/attempts/{attempt_id}/submit", d.QuizHdl().SubmitAttempt)
			r.Post("/attempts/{attempt_id}/answers/{answer_id}/grade", d.QuizHdl().GradeAttemptAnswer)
		})

		// Admin
		r.Route("/api/v1/admin", func(r chi.Router) {
			r.Use(d.MwAuth().JWTMiddleware)

			r.Get("/users", d.AuthHdl().ListUsers)
			r.Patch("/users/{id}/role", d.AuthHdl().ChangeUserRole)
			r.Patch("/users/{id}/ban", d.AuthHdl().SetUserBanned)
		})

		// Static Files (Uploads)
		filesDir := http.Dir("./uploads")
		r.Handle("/static/uploads/*", http.StripPrefix("/static/uploads", http.FileServer(filesDir)))
		r.Handle("/static/*", http.StripPrefix("/static", http.FileServer(filesDir)))

		// Upload
		r.Route("/api/v1/upload", func(r chi.Router) {
			r.Use(d.MwAuth().JWTMiddleware)

			r.Post("/", d.UploadHdl().UploadFile)
		})

		// Categories
		r.Route("/api/v1/categories", func(r chi.Router) {
			r.Get("/", d.CategoryHdl().ListCategories)
			r.Group(func(r chi.Router) {
				r.Use(d.MwAuth().JWTMiddleware)
				r.Post("/", d.CategoryHdl().CreateCategory)
			})
		})

		// Certificates (Public Verification)
		r.Route("/api/v1/certificates", func(r chi.Router) {
			r.Get("/verify/{code}", d.CertHdl().VerifyCertificate)
		})

		// Notifications
		r.Route("/api/v1/notifications", func(r chi.Router) {
			r.Use(d.MwAuth().JWTMiddleware)

			r.Get("/", d.NotifHdl().GetNotifications)
			r.Patch("/{id}/read", d.NotifHdl().MarkAsRead)
			r.Post("/read-all", d.NotifHdl().MarkAllAsRead)
		})

		d.router = r
	}
	return d.router
}
