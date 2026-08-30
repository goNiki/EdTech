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

	authHandler "edtech/internal/interfaces/handlers/auth"
	courseHandler "edtech/internal/interfaces/handlers/courses"
	progressHandler "edtech/internal/interfaces/handlers/progress"
	quizHandler "edtech/internal/interfaces/handlers/quiz"
	mwauth "edtech/internal/interfaces/middleware/auth"
	mwlogger "edtech/internal/interfaces/middleware/logger"

	"edtech/internal/repository"
	authRepo "edtech/internal/repository/auth"
	courseRepo "edtech/internal/repository/course"
	enrolledRepo "edtech/internal/repository/enrollment"
	lessonRepo "edtech/internal/repository/lesson"
	permissionRepo "edtech/internal/repository/permission"
	progressRepo "edtech/internal/repository/progress"
	quizRepo "edtech/internal/repository/quiz"
	refreshRepo "edtech/internal/repository/refresh"
	sectionRepo "edtech/internal/repository/section"

	"edtech/internal/service"
	accessService "edtech/internal/service/access"
	authService "edtech/internal/service/auth"
	courseService "edtech/internal/service/course"
	enrolledService "edtech/internal/service/enrollment"
	lessonService "edtech/internal/service/lesson"
	progressService "edtech/internal/service/progress"
	quizService "edtech/internal/service/quiz"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	// middleware
	mwAuth mwauth.AuthMiddleware
	mwLog  func(http.Handler) http.Handler

	// router
	router http.Handler

	// handlers
	authHdl     *authHandler.AuthHandler
	courseHdl   *courseHandler.CourseHandler
	progressHdl *progressHandler.ProgressHandler
	quizHdl     *quizHandler.QuizHandler

	// services
	accessSvc   service.AccessService
	authSvc     service.AuthService
	courseSvc   service.CourseServices
	enrolledSvc service.EnrolledServices
	lessonSvc   service.LessonServices
	progressSvc service.ProgressServices
	quizSvc     service.QuizServices

	// repositories
	userRepo     repository.UserRepository
	refreshRepo  repository.RefreshRepository
	courseRepo   repository.CourseRepository
	sectionRepo  repository.SectionRepository
	lessonRepo   repository.LessonRepository
	enrolledRepo repository.EnrolledRepository
	permRepo     repository.PermissionsRepository
	progRepo     repository.ProgressRepository
	quizRepo     repository.QuizRepository
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

func (d *diContainer) JWTCfg() config.JWT {
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
		database, err := db.New(d.PostgresCfg())
		if err != nil {
			d.Logger().Error("failed to connect to postgres: " + err.Error())
			os.Exit(1)
		}
		d.db = database
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
		d.jwtManager = jwt.NewJwtManager(d.JWTCfg())
	}
	return d.jwtManager
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

// Services

func (d *diContainer) AccessSvc() service.AccessService {
	if d.accessSvc == nil {
		d.accessSvc = accessService.NewAccessService(d.CourseRepo(), d.EnrolledRepo(), d.PermRepo(), d.DB().Pool)
	}
	return d.accessSvc
}

func (d *diContainer) AuthSvc() service.AuthService {
	if d.authSvc == nil {
		d.authSvc = authService.NewAuthService(d.UserRepo(), d.JWTManager(), d.Hasher(), d.RefreshRepo(), d.DB().Pool, d.TxManager())
	}
	return d.authSvc
}

func (d *diContainer) CourseSvc() service.CourseServices {
	if d.courseSvc == nil {
		d.courseSvc = courseService.NewCourseService(d.CourseRepo(), d.SectionRepo(), d.LessonRepo(), d.AccessSvc(), d.EnrolledRepo(), d.TxManager(), d.DB().Pool)
	}
	return d.courseSvc
}

func (d *diContainer) LessonSvc() service.LessonServices {
	if d.lessonSvc == nil {
		d.lessonSvc = lessonService.NewLessonService(d.CourseRepo(), d.LessonRepo(), d.DB().Pool)
	}
	return d.lessonSvc
}

func (d *diContainer) EnrolledSvc() service.EnrolledServices {
	if d.enrolledSvc == nil {
		d.enrolledSvc = enrolledService.NewEnrolmentService(d.CourseRepo(), d.EnrolledRepo(), d.UserRepo(), d.AccessSvc(), d.DB().Pool, d.TxManager())
	}
	return d.enrolledSvc
}

func (d *diContainer) ProgressSvc() service.ProgressServices {
	if d.progressSvc == nil {
		d.progressSvc = progressService.NewProgressService(d.ProgRepo(), d.LessonRepo(), d.TxManager(), d.DB().Pool)
	}
	return d.progressSvc
}

func (d *diContainer) QuizSvc() service.QuizServices {
	if d.quizSvc == nil {
		d.quizSvc = quizService.NewQuizService(d.QuizRepo(), d.CourseRepo(), d.LessonRepo(), d.AccessSvc(), d.ProgressSvc(), d.TxManager(), d.DB().Pool)
	}
	return d.quizSvc
}

// Handlers

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

// Router

func (d *diContainer) Router() http.Handler {
	if d.router == nil {
		r := chi.NewRouter()

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
			r.Get("/", d.CourseHdl().ListPublicCourses)
			r.Get("/{courseid}", d.CourseHdl().GetCourseByID)
			r.Get("/slug/{slug}", d.CourseHdl().GetCourseBySlug)
			r.Get("/{courseid}/structure", d.CourseHdl().GetCourseStructure)

			r.Group(func(r chi.Router) {
				r.Use(d.MwAuth().JWTMiddleware)

				r.Get("/my", d.CourseHdl().ListMyCourses)
				r.Post("/", d.CourseHdl().CreateCourse)
				r.Patch("/{courseid}", d.CourseHdl().UpdateCourse)
				r.Delete("/{courseid}", d.CourseHdl().DeleteCourse)
				r.Post("/{courseid}/publish", d.CourseHdl().PublishCourse)
				r.Post("/{courseid}/archive", d.CourseHdl().ArchiveCourse)

				r.Get("/{course_id}/progress", d.ProgressHdl().GetCourseProgress)
				r.Get("/{course_id}/progress/lessons", d.ProgressHdl().GetAllLessonProgress)
				r.Get("/{course_id}/quizzes/attempts", d.QuizHdl().ListAttemptsForGrading)
			})
		})

		// Lessons
		r.Route("/api/v1/lessons", func(r chi.Router) {
			r.Use(d.MwAuth().JWTMiddleware)

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

			r.Patch("/users/{id}/role", d.AuthHdl().ChangeUserRole)
			r.Patch("/users/{id}/ban", d.AuthHdl().SetUserBanned)
		})

		d.router = r
	}
	return d.router
}

func (d *diContainer) Close() {
	if d.db != nil {
		d.db.Close()
	}
}
