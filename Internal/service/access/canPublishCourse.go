package access

import (
	"context"
	"edtech/internal/domain"
	errorsAPP "edtech/pkg/errors"
	"errors"
)

func (s *service) CanPublishCourse(ctx context.Context, course *domain.Course, userID int64) (bool, error) {

	role, err := s.enrolledrepo.GetRoleUserInCource(ctx, userID, course.Id)
	if err != nil {
		if errors.Is(err, errorsAPP.ErrUserNotFound) {
			return false, nil
		}
		return false, err
	}

	access, err := s.permissionrepo.HasPermission(ctx, role, domain.ResourceCourse, domain.ActionPublish)
	if err != nil {
		return false, err
	}
	return access, nil

}
