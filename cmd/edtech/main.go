package main

import (
	"edtech/internal/infrastructure/config"
	"edtech/internal/infrastructure/db"
	"edtech/internal/infrastructure/hasher"
	"edtech/internal/infrastructure/jwt"
	"edtech/internal/infrastructure/logger"
	"edtech/internal/infrastructure/logger/sl"
	authHandler "edtech/internal/interfaces/http/handlers/auth"
	mwlogger "edtech/internal/interfaces/middleware"
	authRepo "edtech/internal/repository/auth"
	authUC "edtech/internal/usecase/auth"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// иннициализация логера
	log := logger.Init()

	log.Info("logger init")

	// загрузка конфигурации сервера и БД
	cfg := config.InitConfig()

	log.Info("config init")

	// подключение к базе данных
	postgres, err := db.New(cfg)
	if err != nil {
		log.Error("failed to create connect postgres: ", sl.Error(err))
		return
	}
	defer postgres.Close()

	log.Info("successfully connected to database")

	// инициализация hasher и jwtManager
	hasher := hasher.NewBcryptHasher()
	jwtManager := jwt.NewJwtManager(&cfg.JWTConfig)

	// иннициализация репозитория пользователей
	userrepo := authRepo.NewUserRepo(postgres.Pool)

	// иннициализация сервиса пользователей
	userService := authUC.NewUserUseCase(userrepo, jwtManager, hasher)

	// инициализация маршрутизатора

	authHandler := authHandler.NewAuthHandler(userService)

	// обьявляем роутеры и делаем мидлвееры
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(mwlogger.New(log))

	r.Post("/auth/register", authHandler.Register)

	srv := http.Server{
		Addr:         ":" + cfg.AppConfig.Port,
		Handler:      r,
		ReadTimeout:  cfg.AppConfig.TimeOut,
		IdleTimeout:  cfg.AppConfig.Idletimeout,
		WriteTimeout: 15 * time.Second,
	}

	log.Info("starting server ", slog.String("addr", srv.Addr))
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error("server closed with err", sl.Error(err))
	}
}
