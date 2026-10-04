package course

import (
	"context"
	"fmt"

	"edtech/internal/domain"
	"edtech/pkg/utils"

	"github.com/jackc/pgx/v5"
)

func (s *service) resolveUniqueSlug(ctx context.Context, baseSlug string) (string, error) {
	const op = "service.course.resolveUniqueSlug"

	cleanSlug := utils.NormalizeSlug(baseSlug)
	candidate := cleanSlug
	counter := 1

	for {
		exists, err := s.courserepo.ExistsBySlug(ctx, candidate)
		if err != nil {
			return "", fmt.Errorf("%s: %w", op, err)
		}
		if !exists {
			return candidate, nil
		}
		candidate = fmt.Sprintf("%s-%d", cleanSlug, counter)
		counter++
	}
}

func (s *service) CreateCourse(ctx context.Context, course *domain.Course) (*domain.Course, error) {
	const op = "service.course.CreateCourse"

	// Normalize slug before validation
	course.Slug = utils.NormalizeSlug(course.Slug)

	if err := course.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	course.Status = domain.StatusDraft

	var createdCourse *domain.Course

	err := s.txManager.WithTX(ctx, pgx.TxOptions{}, func(ctx context.Context) error {
		// Resolve unique slug within transaction to prevent uniqueness collisions
		uniqueSlug, err := s.resolveUniqueSlug(ctx, course.Slug)
		if err != nil {
			return err
		}
		course.Slug = uniqueSlug

		var txErr error
		createdCourse, txErr = s.courserepo.CreateCourse(ctx, course)
		if txErr != nil {
			return txErr
		}

		enrollment := domain.EnrolledInCourse{
			UserID:   createdCourse.CreatedBy,
			CourseID: createdCourse.Id,
			Role:     string(domain.CreatorRole),
		}

		if txErr = s.enrolledrepo.EnrollUserToCourse(ctx, enrollment); txErr != nil {
			return txErr
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return createdCourse, nil
}
