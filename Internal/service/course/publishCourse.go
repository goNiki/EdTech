package course

import (
	"context"
	errorsAPP "edtech/pkg/errors"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (s *service) PublishCourse(ctx context.Context, id int64) error {
	course, err := s.repo.GetCourseByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errorsAPP.ErrNotFoundCourse
		}
		return fmt.Errorf("%w: %v", errorsAPP.ErrInternalDB, err)
	}

	//TODO еще дополнительную проверку, чтобы был хотя бы 1 урок на курсе, чтобы его можно было опубликовать
	err = course.Publish()
	if err != nil {
		return fmt.Errorf("%w: %v", errorsAPP.ErrCourseAlredyPublished, err)
	}

	if err := s.repo.PublishCourse(ctx, course); err != nil {
		return fmt.Errorf("%w: %v", errorsAPP.ErrInternalDB, err)
	}

	return nil

}
