package lesson

import (
	"context"
	"edtech/internal/domain"
	errorsAPP "edtech/pkg/errors"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (s *lessonService) ListLessonsByCourseID(ctx context.Context, slug string) (domain.CourseWithLessons, error) {
	const op = "service.lesson.ListLessonsByCourseID"

	course, err := s.courserepo.GetCourseBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, errorsAPP.ErrCourseNotFound) {
			return domain.CourseWithLessons{}, fmt.Errorf("%s: %w", op, errorsAPP.ErrCourseNotFound)
		}
		return domain.CourseWithLessons{}, fmt.Errorf("%s: %w", op, err)
	}

	lessons, err := s.lessonrepo.GetLessonsByCourseID(ctx, course.Id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.CourseWithLessons{}, fmt.Errorf("%s: %w", op, errorsAPP.ErrLessonNotFound)
		}
		return domain.CourseWithLessons{}, fmt.Errorf("%s: %w", op, err)
	}

	lessonsList := domain.CourseWithLessons{
		Course:  *course,
		Lessons: lessons,
	}

	return lessonsList, nil
}
