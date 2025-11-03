package course

import (
	"context"
	"edtech/internal/domain"
	errorsAPP "edtech/pkg/errors"
	"errors"

	"github.com/jackc/pgx/v5"
)

func (s *service) GetCourceByID(ctx context.Context, id int64) (*domain.Course, error) {

	course, err := s.repo.GetCourseByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorsAPP.ErrNotFoundCourse
		}
		return nil, errorsAPP.ErrInternalDB
	}

	return course, nil
}
