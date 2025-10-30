package lesson

import (
	"context"
	"edtech/internal/domain"
	"edtech/internal/infrastructure/logger"
	"edtech/internal/infrastructure/logger/sl"
	errorsAPP "edtech/pkg/errors"
	"errors"

	"github.com/jackc/pgx/v5"
)

func (s *service) ListLessonsByCourseID(ctx context.Context, slug string) (domain.CourseWithLessons, error) {

	const op = "usecase.course.listlessonsbycourseid"

	log := logger.GetLogger(ctx, op)

	course, err := s.courserepo.GetCourseBySlug(ctx, slug)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error("course is not found", sl.Error(err))
			return domain.CourseWithLessons{}, errorsAPP.ErrNotFoundCourse
		}
		log.Error("DB error", sl.Error(err))
		return domain.CourseWithLessons{}, errorsAPP.ErrInternalDB
	}

	lessons, err := s.lessonrepo.GetLessonsByCourseID(ctx, course.Id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error("lessons are not found", sl.Error(err))
			return domain.CourseWithLessons{}, errorsAPP.ErrNotFoundLesson
		}
		log.Error("DB error", sl.Error(err))
		return domain.CourseWithLessons{}, errorsAPP.ErrInternalDB
	}

	lessonsList := domain.CourseWithLessons{
		Course:  *course,
		Lessons: lessons,
	}

	return lessonsList, nil
}
