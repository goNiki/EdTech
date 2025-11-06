package course

import (
	"context"
	"edtech/internal/domain"
	errorsAPP "edtech/pkg/errors"
	"edtech/pkg/utils"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (s *service) CreateCourse(ctx context.Context, course *domain.Course) (int64, error) {

	if err := utils.ValidateCourse(course.Title, course.Slug, course.CreatedBy, course.Visibility, course.Status); err != nil {
		return 0, err
	}

	if _, err := s.courserepo.GetCourseBySlug(ctx, course.Slug); !errors.Is(err, pgx.ErrNoRows) {
		if err == nil {
			return 0, errorsAPP.ErrSlugAlreadyExists
		}
		return 0, fmt.Errorf("%w: %v", errorsAPP.ErrInternalDB, err)
	}

	course.Status = "draft"

	courseID, err := s.courserepo.CreateCourse(ctx, course)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", errorsAPP.ErrInternalDB, err)
	}

	return courseID, nil

}
