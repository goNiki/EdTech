package analytics

import (
	"context"
	"fmt"

	"edtech/internal/domain"
)

func (s *analyticsService) GetCourseAnalytics(ctx context.Context, teacherID, courseID int64) (domain.CourseAnalyticsSummary, error) {
	const op = "service.analytics.GetCourseAnalytics"

	if err := s.checkTeacherAccess(ctx, teacherID, courseID); err != nil {
		return domain.CourseAnalyticsSummary{}, fmt.Errorf("%s: %w", op, err)
	}

	summary, err := s.analyticsRepo.GetCourseAnalyticsSummary(ctx, courseID)
	if err != nil {
		return domain.CourseAnalyticsSummary{}, fmt.Errorf("%s: %w", op, err)
	}

	return summary, nil
}
