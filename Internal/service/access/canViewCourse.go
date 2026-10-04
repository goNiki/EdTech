package access

import (
	"context"
	"edtech/internal/domain"
	errorsAPP "edtech/pkg/errors"
	"errors"
)

func (s *service) CanViewCourse(ctx context.Context, course *domain.Course, userID int64) (bool, error) {

	if course.Visibility == domain.VisibilityPublic && course.Status == domain.StatusPublished {
		return true, nil
	}

	if course.CreatedBy == userID {
		return true, nil
	}

	role, err := s.enrolledrepo.GetRoleUserInCourse(ctx, userID, course.Id)
	if err != nil {
		if errors.Is(err, errorsAPP.ErrUserNotFound) {
			return false, nil
		}
		return false, err
	}

	switch role {
	case string(domain.StudentRole):
		if course.Status == domain.StatusPublished {
			return true, nil
		}
	case string(domain.TeacherRole):
		return true, nil
	case string(domain.CreatorRole):
		return true, nil
	}
	return false, nil
}
