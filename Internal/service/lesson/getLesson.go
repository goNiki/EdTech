package lesson

import (
	"context"
	"fmt"

	"edtech/internal/domain"
	"edtech/internal/service/quiz"
)

func (s *lessonService) GetLesson(ctx context.Context, userID int64, id int64) (*domain.Lesson, error) {
	const op = "service.lesson.GetLesson"

	lesson, err := s.lessonrepo.GetLessonByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	canEdit := false
	if userID > 0 && s.courserepo != nil && s.accessService != nil {
		course, cErr := s.courserepo.GetCourseByID(ctx, lesson.CourseID)
		if cErr == nil && course != nil {
			if ok, aErr := s.accessService.CanEditCourse(ctx, course, userID); aErr == nil {
				canEdit = ok
			}
		}
	}

	// If not course editor/admin, sanitize quizzes to prevent cheating
	if !canEdit && lesson.Content != "" {
		sanitized, sErr := quiz.SanitizeLessonContentForStudent(lesson.Content)
		if sErr == nil {
			lesson.Content = sanitized
		}
	}

	return lesson, nil
}

