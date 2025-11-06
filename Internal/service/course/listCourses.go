package course

import (
	"context"
	"edtech/internal/domain"
	errorsAPP "edtech/pkg/errors"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (s *service) ListCourses(ctx context.Context, page int64, pageSize int64) (domain.PaginatedCourses, error) {

	offset := (page - 1) * pageSize

	corses, err := s.courserepo.ListCourses(ctx, pageSize, offset)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.PaginatedCourses{}, errorsAPP.ErrNotFoundCourse
		}
		return domain.PaginatedCourses{}, fmt.Errorf("%w: %c", errorsAPP.ErrInternalDB, err)
	}

	total, err := s.courserepo.CountCourse(ctx)
	if err != nil {
		return domain.PaginatedCourses{}, fmt.Errorf("%w: %c", errorsAPP.ErrInternalDB, err)
	}
	coursesList := domain.PaginatedCourses{
		Courses:  corses,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}

	return coursesList, nil
}
