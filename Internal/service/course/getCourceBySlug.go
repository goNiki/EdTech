package course

import (
	"context"
	"edtech/internal/domain"
	errorsAPP "edtech/pkg/errors"
	"errors"

	"github.com/jackc/pgx/v5"
)

func (s *service) GetCourseBySlug(ctx context.Context, slug string) (*domain.Course, error) {

	course, err := s.repo.GetCourseBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorsAPP.ErrNotFoundCourse
		}
		return nil, errorsAPP.ErrInternalDB
	}

	return course, nil
}
