package analytics

import (
	"context"
	"fmt"

	"edtech/internal/domain"
)

func (s *analyticsService) ListPendingHomeworks(ctx context.Context, teacherID, courseID int64, page, pageSize int64) ([]domain.PendingHomeworkItem, int64, error) {
	const op = "service.analytics.ListPendingHomeworks"

	if err := s.checkTeacherAccess(ctx, teacherID, courseID); err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	items, total, err := s.analyticsRepo.ListPendingHomeworks(ctx, courseID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	return items, total, nil
}
