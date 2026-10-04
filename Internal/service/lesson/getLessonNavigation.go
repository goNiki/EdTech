package lesson

import (
	"context"
	"fmt"
	"sort"

	"edtech/internal/domain"
	errorsAPP "edtech/pkg/errors"
)

func (s *lessonService) GetLessonNavigationContext(ctx context.Context, userID int64, lessonID int64) (*domain.LessonNavigationContext, error) {
	const op = "service.lesson.GetLessonNavigationContext"

	// 1. Найти текущий урок
	lesson, err := s.lessonrepo.GetLessonByID(ctx, lessonID)
	if err != nil {
		return nil, fmt.Errorf("%s: get current lesson: %w", op, err)
	}

	// 2. Найти курс
	course, err := s.courserepo.GetCourseByID(ctx, lesson.CourseID)
	if err != nil {
		return nil, fmt.Errorf("%s: get course: %w", op, err)
	}

	// 3. Проверить доступ пользователя к курсу
	if s.accessService != nil {
		canView, aErr := s.accessService.CanViewCourse(ctx, course, userID)
		if aErr != nil {
			return nil, fmt.Errorf("%s: check course access: %w", op, aErr)
		}
		if !canView {
			return nil, fmt.Errorf("%s: %w", op, errorsAPP.ErrForbidden)
		}
	}

	// 4. Получить секции курса (отсортированы по position ASC)
	sections, err := s.sectionrepo.ListSectionsByCourseID(ctx, course.Id)
	if err != nil {
		return nil, fmt.Errorf("%s: list sections: %w", op, err)
	}

	// 5. Получить уроки курса
	lessons, err := s.lessonrepo.GetLessonsByCourseID(ctx, course.Id)
	if err != nil {
		return nil, fmt.Errorf("%s: get lessons: %w", op, err)
	}

	// 6. Получить прогресс студента по курсу (если передан валидный userID)
	progressMap := make(map[int64]*domain.LessonProgress)
	if s.progressrepo != nil && userID > 0 {
		progs, pErr := s.progressrepo.GetAllLessonProgressByCourse(ctx, userID, course.Id)
		if pErr == nil {
			for i := range progs {
				progressMap[progs[i].LessonID] = &progs[i]
			}
		}
	}

	// 7. Сгруппировать уроки по секциям
	lessonsBySection := make(map[int64][]domain.Lesson)
	var unsectioned []domain.Lesson

	for _, l := range lessons {
		if l.SectionID != nil {
			lessonsBySection[*l.SectionID] = append(lessonsBySection[*l.SectionID], l)
		} else {
			unsectioned = append(unsectioned, l)
		}
	}

	// 8. Построить syllabus и плоский упорядоченный список flatLessons
	var flatLessons []domain.Lesson
	syllabus := make([]domain.LessonNavSection, 0, len(sections))

	for _, sec := range sections {
		secLessons := lessonsBySection[sec.ID]
		sort.SliceStable(secLessons, func(i, j int) bool {
			return secLessons[i].Position < secLessons[j].Position
		})

		items := make([]domain.LessonNavItem, 0, len(secLessons))
		for _, l := range secLessons {
			var isCompleted bool
			score := 0
			if prog, ok := progressMap[l.ID]; ok {
				isCompleted = prog.Status == domain.ProgressStatusCompleted
				if prog.Score != nil {
					score = *prog.Score
				}
			}
			items = append(items, domain.LessonNavItem{
				ID:          l.ID,
				Title:       l.Title,
				Position:    l.Position,
				IsCompleted: isCompleted,
				Score:       score,
			})
		}

		syllabus = append(syllabus, domain.LessonNavSection{
			SectionID:    sec.ID,
			SectionTitle: sec.Title,
			Position:     sec.Position,
			Lessons:      items,
		})

		flatLessons = append(flatLessons, secLessons...)
	}

	if len(unsectioned) > 0 {
		sort.SliceStable(unsectioned, func(i, j int) bool {
			return unsectioned[i].Position < unsectioned[j].Position
		})

		items := make([]domain.LessonNavItem, 0, len(unsectioned))
		for _, l := range unsectioned {
			var isCompleted bool
			score := 0
			if prog, ok := progressMap[l.ID]; ok {
				isCompleted = prog.Status == domain.ProgressStatusCompleted
				if prog.Score != nil {
					score = *prog.Score
				}
			}
			items = append(items, domain.LessonNavItem{
				ID:          l.ID,
				Title:       l.Title,
				Position:    l.Position,
				IsCompleted: isCompleted,
				Score:       score,
			})
		}

		syllabus = append(syllabus, domain.LessonNavSection{
			SectionID:    0,
			SectionTitle: "Общие уроки",
			Position:     len(sections) + 1,
			Lessons:      items,
		})

		flatLessons = append(flatLessons, unsectioned...)
	}

	// 9. Найти текущий урок в flatLessons для определения prev и next
	currIdx := -1
	for i, l := range flatLessons {
		if l.ID == lessonID {
			currIdx = i
			break
		}
	}

	var prevLesson *domain.LessonNavNeighbor
	if currIdx > 0 {
		prevLesson = &domain.LessonNavNeighbor{
			ID:    flatLessons[currIdx-1].ID,
			Title: flatLessons[currIdx-1].Title,
		}
	}

	var nextLesson *domain.LessonNavNeighbor
	if currIdx >= 0 && currIdx < len(flatLessons)-1 {
		nextLesson = &domain.LessonNavNeighbor{
			ID:    flatLessons[currIdx+1].ID,
			Title: flatLessons[currIdx+1].Title,
		}
	}

	return &domain.LessonNavigationContext{
		CurrentLesson: domain.LessonNavCurrent{
			ID:        lesson.ID,
			Title:     lesson.Title,
			Position:  lesson.Position,
			SectionID: lesson.SectionID,
		},
		Course: domain.LessonNavCourse{
			ID:    course.Id,
			Title: course.Title,
			Slug:  course.Slug,
		},
		PrevLesson: prevLesson,
		NextLesson: nextLesson,
		Syllabus:   syllabus,
	}, nil
}
