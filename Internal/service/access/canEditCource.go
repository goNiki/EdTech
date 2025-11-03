package access

import (
	"context"
	"edtech/internal/domain"
	errorsAPP "edtech/pkg/errors"
	"errors"
)

func (s *service) CanEditCourse(ctx context.Context, course *domain.Course, userID int64) (bool, error) {

	if userID == course.CreatedBy {
		return true, nil
	}

	role, err := s.enrolledrepo.GetRoleUserInCource(ctx, userID, course.Id)
	if err != nil {
		if errors.Is(err, errorsAPP.ErrUserNotFound) {
			return false, nil
		}
		return false, err
	}

	access, err := s.permissionrepo.HasPermission(ctx, role, domain.ResourceCource, domain.ActionEdit)
	if err != nil {
		return false, err
	}

	return access, nil

}
