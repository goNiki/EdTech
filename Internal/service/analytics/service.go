package analytics

import (
	"context"
	"fmt"

	"edtech/internal/infrastructure/db"
	"edtech/internal/infrastructure/txmanager"
	"edtech/internal/repository"
	"edtech/internal/service"
	errorsAPP "edtech/pkg/errors"
)

type analyticsService struct {
	analyticsRepo repository.AnalyticsRepository
	courseRepo    repository.CourseRepository
	accessService service.AccessService
	txManager     txmanager.TransactionManager
	db            db.QueryExecutor
}

func NewAnalyticsService(
	analyticsRepo repository.AnalyticsRepository,
	courseRepo repository.CourseRepository,
	accessService service.AccessService,
	txManager txmanager.TransactionManager,
	database db.QueryExecutor,
) service.AnalyticsServices {
	return &analyticsService{
		analyticsRepo: analyticsRepo,
		courseRepo:    courseRepo,
		accessService: accessService,
		txManager:     txManager,
		db:            database,
	}
}

func (s *analyticsService) checkTeacherAccess(ctx context.Context, teacherID, courseID int64) error {
	course, err := s.courseRepo.GetCourseByID(ctx, s.db, courseID)
	if err != nil {
		return fmt.Errorf("get course: %w", err)
	}

	canEdit, err := s.accessService.CanEditCourse(ctx, course, teacherID)
	if err != nil {
		return fmt.Errorf("%w: %w", errorsAPP.ErrCheckingPermissions, err)
	}
	if !canEdit {
		return errorsAPP.ErrForbidden
	}

	return nil
}
