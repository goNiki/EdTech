package access

import (
	"context"

	"edtech/internal/domain"
)

func (s *service) BuildCoursePermissions(ctx context.Context, course *domain.Course, userID int64) (domain.CoursePermissions, error) {
	if course == nil {
		return domain.CoursePermissions{}, nil
	}

	isCreator := course.CreatedBy == userID

	// Получаем роль пользователя в данном курсе за 1 запрос к БД
	role, err := s.enrolledrepo.GetRoleUserInCourse(ctx, s.db, userID, course.Id)
	isEnrolled := err == nil && role != ""

	// 1. CanView
	canView := false
	if (course.Visibility == domain.VisibilityPublic && course.Status == domain.StatusPublished) || isCreator {
		canView = true
	} else if isEnrolled {
		if role == string(domain.StudentRole) && course.Status == domain.StatusPublished {
			canView = true
		} else if role == string(domain.TeacherRole) || role == string(domain.CreatorRole) {
			canView = true
		}
	}

	// 2. CanEdit
	canEdit := false
	if isCreator {
		canEdit = true
	} else if isEnrolled && (role == string(domain.TeacherRole) || role == string(domain.CreatorRole)) {
		canEdit = true
	}

	// 3. CanDelete
	canDelete := isCreator

	// 4. CanPublish
	canPublish := isCreator

	// 5. CanEnroll (самозапись)
	canEnroll := false
	if course.Status == domain.StatusPublished && course.Visibility == domain.VisibilityPublic && !isEnrolled && !isCreator {
		canEnroll = true
	}

	// 6. CanManageUsers
	canManageUsers := isCreator

	return domain.CoursePermissions{
		CanView:        canView,
		CanEdit:        canEdit,
		CanDelete:      canDelete,
		CanPublish:     canPublish,
		CanEnroll:      canEnroll,
		CanManageUsers: canManageUsers,
	}, nil
}
