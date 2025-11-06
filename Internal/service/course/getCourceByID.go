package course

import (
	"context"
	"edtech/internal/domain"
	errorsAPP "edtech/pkg/errors"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (s *service) GetCourseByID(ctx context.Context, id int64) (*domain.Course, error) {

	course, err := s.courserepo.GetCourseByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorsAPP.ErrNotFoundCourse
		}
		return nil, fmt.Errorf("%w: %v", errorsAPP.ErrInternalDB, err)
	}

	return course, nil
}
