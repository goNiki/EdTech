package main

import (
	"edtech/Internal/infrastructure/config"
	"edtech/Internal/infrastructure/db"
	"edtech/Internal/infrastructure/logger"
	"edtech/Internal/infrastructure/logger/sl"
	mwlogger "edtech/Internal/interfaces/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// иннициализация логера
	log := logger.Init()

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

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(mwlogger.New(log))

}
