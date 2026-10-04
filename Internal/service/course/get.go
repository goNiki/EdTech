package course

import (
	"context"
	"edtech/internal/domain"
	"fmt"
)

func (s *service) GetCourseByID(ctx context.Context, id int64) (*domain.Course, error) {
	const op = "service.course.GetCourseByID"

	course, err := s.courserepo.GetCourseByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return course, nil
}

func (s *service) GetCourseBySlug(ctx context.Context, slug string) (*domain.Course, error) {
	const op = "service.course.GetCourseBySlug"

	course, err := s.courserepo.GetCourseBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return course, nil
}

func (s *service) GetCourseWithLessons(ctx context.Context, courseID int64) (domain.CourseWithLessons, error) {
	const op = "service.course.GetCourseWithLessons"

	course, err := s.courserepo.GetCourseByID(ctx, courseID)
	if err != nil {
		return domain.CourseWithLessons{}, fmt.Errorf("%s: %w", op, err)
	}

	// 2. Получаем уроки курса
	lessons, err := s.lessonrepo.GetLessonsByCourseID(ctx, courseID)
	if err != nil {
		return domain.CourseWithLessons{}, fmt.Errorf("%s: %w", op, err)
	}

	if lessons == nil {
		lessons = []domain.Lesson{}
	}

	return domain.CourseWithLessons{
		Course:  *course,
		Lessons: lessons,
	}, nil
}
