package access

import (
	"context"
	"edtech/internal/domain"
	errorsAPP "edtech/pkg/errors"
)

func (s *service) CanAccessCourseObject(ctx context.Context, course *domain.Course, userID int64, action string) (bool, error) {
	switch action {
	case domain.ActionView:
		return s.CanViewCourse(ctx, course, userID)
	case domain.ActionEdit:
		return s.CanEditCourse(ctx, course, userID)
	case domain.ActionDelete:
		return s.CanDeleteCourse(ctx, course, userID)
	case domain.ActionEnroll:
		return s.CanSelfEnrollCourse(ctx, course, userID)
	case domain.ActionManageusers:
		return s.CanManageCourseUsers(ctx, course, userID)
	case domain.ActionPublish:
		return s.CanPublishCourse(ctx, course, userID)
	}
	return false, errorsAPP.ErrInvalidAction

}
