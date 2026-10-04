package analytics

import (
	"context"
	"fmt"

	"edtech/internal/domain"
)

func (s *analyticsService) ListTeacherPendingHomeworks(
	ctx context.Context,
	teacherID int64,
	courseID int64,
	page, pageSize int64,
) (*domain.TeacherPendingHomeworksResult, error) {
	const op = "service.analytics.ListTeacherPendingHomeworks"

	if courseID > 0 {
		if err := s.checkTeacherAccess(ctx, teacherID, courseID); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	items, total, summary, err := s.analyticsRepo.ListTeacherPendingHomeworks(ctx, teacherID, courseID, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if items == nil {
		items = []domain.PendingHomeworkItem{}
	}
	if summary == nil {
		summary = []domain.CoursePendingSummaryItem{}
	}

	return &domain.TeacherPendingHomeworksResult{
		Items:          items,
		Total:          total,
		CoursesSummary: summary,
	}, nil
}
