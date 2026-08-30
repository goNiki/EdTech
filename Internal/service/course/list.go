package course

import (
	"context"
	"fmt"

	"edtech/internal/domain"
)

func (s *service) ListMyCourses(ctx context.Context, input *domain.InputListMyCourse) (domain.PaginatedCourses, error) {
	const op = "service.course.ListMyCourses"

	if input == nil {
		input = &domain.InputListMyCourse{}
	}

	// 1. Валидация и санитизация пагинации
	input.Pagination.Sanitize()

	// 2. Подсчет общего количества с учетом роли и фильтров
	total, err := s.courserepo.CountEnrolledCourses(ctx, s.db, input)
	if err != nil {
		return domain.PaginatedCourses{}, fmt.Errorf("%s: count enrolled courses: %w", op, err)
	}

	if total == 0 {
		return domain.PaginatedCourses{
			Courses:  []domain.Course{},
			Page:     input.Pagination.Page,
			PageSize: input.Pagination.PageSize,
			Total:    0,
		}, nil
	}

	// 3. Получение отфильтрованного и отсортированного списка курсов
	courses, err := s.courserepo.ListEnrolledCourses(ctx, s.db, input)
	if err != nil {
		return domain.PaginatedCourses{}, fmt.Errorf("%s: list enrolled courses: %w", op, err)
	}

	return domain.PaginatedCourses{
		Courses:  courses,
		Page:     input.Pagination.Page,
		PageSize: input.Pagination.PageSize,
		Total:    total,
	}, nil
}

func (s *service) ListPublicCourses(ctx context.Context, pagination domain.Pagination, filter domain.CourseFilter) (domain.PaginatedCourses, error) {
	const op = "service.course.ListPublicCourses"

	// 1. Валидация и санитизация пагинации
	pagination.Sanitize()

	// 2. Получаем общее количество опубликованных публичных курсов с учетом фильтров
	total, err := s.courserepo.CountCourses(ctx, s.db, filter)
	if err != nil {
		return domain.PaginatedCourses{}, fmt.Errorf("%s: count courses: %w", op, err)
	}

	// 3. Если курсов 0 — сразу возвращаем пустую структуру без лишнего запроса
	if total == 0 {
		return domain.PaginatedCourses{
			Courses:  []domain.Course{},
			Page:     pagination.Page,
			PageSize: pagination.PageSize,
			Total:    0,
		}, nil
	}

	// 4. Получаем страницу отфильтрованных и отсортированных курсов
	courses, err := s.courserepo.ListPublicCourses(ctx, s.db, pagination, filter)
	if err != nil {
		return domain.PaginatedCourses{}, fmt.Errorf("%s: list courses: %w", op, err)
	}

	return domain.PaginatedCourses{
		Courses:  courses,
		Page:     pagination.Page,
		PageSize: pagination.PageSize,
		Total:    total,
	}, nil
}
