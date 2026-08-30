package access

import (
	"context"
	"edtech/internal/domain"
	errorsAPP "edtech/pkg/errors"
	"errors"
)

func (s *service) CanManageCourseUsers(ctx context.Context, course *domain.Course, userID int64) (bool, error) {

	if userID == course.CreatedBy {
		return true, nil
	}

	role, err := s.enrolledrepo.GetRoleUserInCourse(ctx, s.db, userID, course.Id)
	if err != nil {
		if errors.Is(err, errorsAPP.ErrUserNotFound) {
			return false, nil
		}
		return false, err
	}

	access, err := s.permissionrepo.HasPermission(ctx, s.db, role, domain.ResourceCourse, domain.ActionManageusers)
	if err != nil {
		return false, err
	}

	return access, nil
}
