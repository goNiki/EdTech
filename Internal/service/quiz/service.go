package quiz

import (
	"edtech/internal/infrastructure/txmanager"
	"edtech/internal/repository"
	services "edtech/internal/service"
)

type service struct {
	quizRepo        repository.QuizRepository
	courseRepo      repository.CourseRepository
	lessonRepo      repository.LessonRepository
	accessService   services.AccessService
	progressService services.ProgressServices
	txManager       txmanager.TransactionManager
	notificationSvc services.NotificationServices
}

func NewQuizService(
	quizRepo repository.QuizRepository,
	courseRepo repository.CourseRepository,
	lessonRepo repository.LessonRepository,
	accessService services.AccessService,
	progressService services.ProgressServices,
	txManager txmanager.TransactionManager,
	notificationSvc services.NotificationServices,
) services.QuizServices {
	return &service{
		quizRepo:        quizRepo,
		courseRepo:      courseRepo,
		lessonRepo:      lessonRepo,
		accessService:   accessService,
		progressService: progressService,
		txManager:       txManager,
		notificationSvc: notificationSvc,
	}
}
