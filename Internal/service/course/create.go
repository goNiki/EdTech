package course

import (
	"context"
	"fmt"

	"edtech/internal/domain"
	"edtech/internal/infrastructure/db"

	"github.com/jackc/pgx/v5"
)

func (s *service) CreateCourse(ctx context.Context, course *domain.Course) (*domain.Course, error) {
	const op = "service.course.CreateCourse"

	if err := course.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	course.Status = domain.StatusDraft

	var createdCourse *domain.Course

	err := s.txManager.WithTX(ctx, pgx.TxOptions{}, func(ctx context.Context, q db.QueryExecutor) error {
		var txErr error
		createdCourse, txErr = s.courserepo.CreateCourse(ctx, q, course)
		if txErr != nil {
			return txErr
		}

		enrollment := domain.EnrolledInCourse{
			UserID:   createdCourse.CreatedBy,
			CourseID: createdCourse.Id,
			Role:     string(domain.CreatorRole),
		}

		if txErr = s.enrolledrepo.EnrollUserToCourse(ctx, q, enrollment); txErr != nil {
			return txErr
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return createdCourse, nil
}
