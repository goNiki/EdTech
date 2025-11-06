package courses

import (
	"context"
	"edtech/internal/domain"
	"edtech/internal/dto"
)

// TODO сделать данную функцию слое access, только сделать корректную конвертацию структур domain в dto
func (h *handler) BuildCoursePermissions(ctx context.Context, course *domain.Course, userID int64) dto.CoursePermissions {
	canView, _ := h.accessService.CanViewCourse(ctx, course, userID)
	canEdit, _ := h.accessService.CanEditCourse(ctx, course, userID)
	canDelete, _ := h.accessService.CanDeleteCourse(ctx, course, userID)
	canPublish, _ := h.accessService.CanPublishCourse(ctx, course, userID)
	canEnroll, _ := h.accessService.CanSelfEnrollCourse(ctx, course, userID)
	canManageCourse, _ := h.accessService.CanManageCourseUsers(ctx, course, userID)

	return dto.CoursePermissions{
		CanView:        canView,
		CanEdit:        canEdit,
		CanDelete:      canDelete,
		CanPublish:     canPublish,
		CanEnroll:      canEnroll,
		CanManageUsers: canManageCourse,
	}
}
