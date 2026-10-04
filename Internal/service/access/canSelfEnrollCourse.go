package access

import (
	"context"
	"edtech/internal/domain"
)

func (s *service) CanSelfEnrollCourse(ctx context.Context, course *domain.Course, userID int64) (bool, error) {

	exists, err := s.enrolledrepo.UserExistCourse(ctx, userID, course.Id)
	if err != nil {
		return false, err
	}

	if exists {
		return false, nil
	}

	return true, nil

}
