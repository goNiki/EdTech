package course

import (
	"context"
	"edtech/internal/domain"
	errorsAPP "edtech/pkg/errors"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (s *service) GetCourseBySlug(ctx context.Context, slug string) (*domain.Course, error) {

	course, err := s.courserepo.GetCourseBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorsAPP.ErrNotFoundCourse
		}
		return nil, fmt.Errorf("%w: %v", errorsAPP.ErrInternalDB, err)
	}

	return course, nil
}
