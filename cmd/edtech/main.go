package main

import (
	"edtech/internal/infrastructure/config"
	"edtech/internal/infrastructure/db"
	"edtech/internal/infrastructure/hasher"
	"edtech/internal/infrastructure/jwt"
	"edtech/internal/infrastructure/logger"
	"edtech/internal/infrastructure/logger/sl"
	"edtech/internal/infrastructure/txmanager"
	authHandler "edtech/internal/interfaces/handlers/auth"
	mwauth "edtech/internal/interfaces/middleware/auth"
	mwlogger "edtech/internal/interfaces/middleware/logger"
	authRepo "edtech/internal/repository/auth"
	refreshRepo "edtech/internal/repository/refresh"
	authService "edtech/internal/service/auth"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// инициализация логера
	log := logger.Init()
	log.Info("logger initialized")

	// загрузка конфигурации сервера и БД
	cfg, err := config.Load(".env")
	if err != nil {
		log.Error("failed to load config: ", sl.Error(err))
		return
	}
	log.Info("config loaded")

	// подключение к базе данных
	postgres, err := db.New(cfg.Postgres)
	if err != nil {
		log.Error("failed to connect to postgres: ", sl.Error(err))
		return
	}
	defer postgres.Close()
	log.Info("successfully connected to database")

	// инициализация txmanager, hasher и jwtManager
	txManager := txmanager.NewTxManager(postgres)
	hasher := hasher.NewHasher()
	jwtManager := jwt.NewJwtManager(cfg.JWT)

	// инициализация middleware
	mwAuth := mwauth.NewAuthMiddleware(jwtManager)
	mwLog := mwlogger.NewLoggerMiddleware(log)

	// инициализация репозиториев
	userRepo := authRepo.NewAuthRepo(postgres.Pool)
	refreshRepo := refreshRepo.NewRefreshRepo(postgres.Pool)

	// инициализация сервисов
	authSvc := authService.NewAuthService(userRepo, jwtManager, hasher, refreshRepo, postgres.Pool, txManager)

	// инициализация хендлеров
	authHdl := authHandler.NewAuthHandler(authSvc, mwAuth)

	// настройка маршрутизатора
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(mwLog.Log)
	r.Use(middleware.Recoverer)

	// Публичные эндпоинты авторизации
	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/register", authHdl.Register)
		r.Post("/login", authHdl.Login)
		r.Post("/refresh", authHdl.RefreshToken)
		r.Post("/logout", authHdl.Logout)

		// Защищенные эндпоинты профиля
		r.Group(func(r chi.Router) {
			r.Use(mwAuth.JWTMiddleware)

			r.Get("/me", authHdl.GetCurrentUser)
			r.Patch("/profile", authHdl.UpdateProfile)
			r.Post("/change-password", authHdl.ChangePassword)
			r.Post("/verify-email", authHdl.VerifyEmail)
		})
	})

	// Защищенные эндпоинты администратора
	r.Route("/api/v1/admin", func(r chi.Router) {
		r.Use(mwAuth.JWTMiddleware)

		r.Patch("/users/{id}/role", authHdl.ChangeUserRole)
		r.Patch("/users/{id}/ban", authHdl.SetUserBanned)
	})

	srv := http.Server{
		Addr:         ":" + cfg.Server.Port(),
		Handler:      r,
		ReadTimeout:  cfg.Server.TimeOut(),
		IdleTimeout:  cfg.Server.Idletimeout(),
		WriteTimeout: 15 * time.Second,
	}

	log.Info("starting server", slog.String("addr", srv.Addr))
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error("server closed with error", sl.Error(err))
	}
}
