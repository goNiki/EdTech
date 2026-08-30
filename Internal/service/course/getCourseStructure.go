package course

import (
	"context"
	"fmt"

	"edtech/internal/domain"
)

func (s *service) GetCourseStructure(ctx context.Context, courseID int64) (domain.CourseStructure, error) {
	const op = "service.course.GetCourseStructure"

	// 1. Получаем курс
	course, err := s.courserepo.GetCourseByID(ctx, s.db, courseID)
	if err != nil {
		return domain.CourseStructure{}, fmt.Errorf("%s: get course: %w", op, err)
	}

	// 2. Получаем секции курса
	sections, err := s.sectionrepo.ListSectionsByCourseID(ctx, s.db, courseID)
	if err != nil {
		return domain.CourseStructure{}, fmt.Errorf("%s: get sections: %w", op, err)
	}

	// 3. Получаем уроки курса
	lessons, err := s.lessonrepo.GetLessonsByCourseID(ctx, s.db, courseID)
	if err != nil {
		return domain.CourseStructure{}, fmt.Errorf("%s: get lessons: %w", op, err)
	}

	// 4. Группируем уроки по секциям
	lessonsBySection := make(map[int64][]domain.Lesson)
	var unsectioned []domain.Lesson

	for _, lesson := range lessons {
		if lesson.SectionID != nil {
			lessonsBySection[*lesson.SectionID] = append(lessonsBySection[*lesson.SectionID], lesson)
		} else {
			unsectioned = append(unsectioned, lesson)
		}
	}

	// 5. Собираем итоговые секции с уроками
	sectionsWithLessons := make([]domain.SectionWithLessons, 0, len(sections))
	for _, sec := range sections {
		secLessons := lessonsBySection[sec.ID]
		if secLessons == nil {
			secLessons = []domain.Lesson{}
		}
		sectionsWithLessons = append(sectionsWithLessons, domain.SectionWithLessons{
			Section: sec,
			Lessons: secLessons,
		})
	}

	if unsectioned == nil {
		unsectioned = []domain.Lesson{}
	}

	return domain.CourseStructure{
		Course:             *course,
		Sections:           sectionsWithLessons,
		UnsectionedLessons: unsectioned,
	}, nil
}
