package course

import (
	"context"
	"edtech/internal/domain"
	"edtech/internal/infrastructure/logger"
	"edtech/internal/infrastructure/logger/sl"
	errorsAPP "edtech/pkg/errors"
	"errors"

	"github.com/jackc/pgx/v5"
)

func (s *service) ListCourses(ctx context.Context, page int, pageSize int) (domain.PaginatedCourses, error) {

	const op = "usecase.course.listCourses"

	log := logger.GetLogger(ctx, op)

	offset := (page - 1) * pageSize

	corses, err := s.repo.ListCourses(ctx, pageSize, offset)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error("courses is not found", sl.Error(err))
			return domain.PaginatedCourses{}, errorsAPP.ErrNotFoundCourse
		}
		log.Error("DB Error: ", sl.Error(err))
		return domain.PaginatedCourses{}, errorsAPP.ErrInternalDB
	}

	total, err := s.repo.CountCourse(ctx)
	if err != nil {
		log.Error("DB Error: ", sl.Error(err))
		return domain.PaginatedCourses{}, errorsAPP.ErrInternalDB
	}
	coursesList := domain.PaginatedCourses{
		Courses:  corses,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}

	return coursesList, nil
}
