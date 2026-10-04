package analytics

import (
	"context"
	"fmt"

	"edtech/internal/domain"
)

func (s *analyticsService) GetStudentDrilldown(ctx context.Context, teacherID, courseID, studentID int64) (*domain.StudentDrilldownReport, error) {
	const op = "service.analytics.GetStudentDrilldown"

	if err := s.checkTeacherAccess(ctx, teacherID, courseID); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	report, err := s.analyticsRepo.GetStudentDrilldown(ctx, studentID, courseID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return report, nil
}
