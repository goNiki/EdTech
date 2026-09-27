package progress

import (
	"edtech/internal/infrastructure/db"
	"edtech/internal/infrastructure/txmanager"
	"edtech/internal/repository"
	services "edtech/internal/service"
)

var _ services.ProgressServices = (*service)(nil)

type service struct {
	progressRepo repository.ProgressRepository
	lessonRepo   repository.LessonRepository
	quizRepo     repository.QuizRepository
	txManager    txmanager.TransactionManager
	db           db.QueryExecutor
}

func NewProgressService(
	progressRepo repository.ProgressRepository,
	lessonRepo repository.LessonRepository,
	quizRepo repository.QuizRepository,
	txManager txmanager.TransactionManager,
	database db.QueryExecutor,
) *service {
	return &service{
		progressRepo: progressRepo,
		lessonRepo:   lessonRepo,
		quizRepo:     quizRepo,
		txManager:    txManager,
		db:           database,
	}
}
