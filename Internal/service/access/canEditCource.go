package access

import (
	"context"
	"edtech/internal/domain"
	errorsAPP "edtech/pkg/errors"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (s *service) CanEditCourse(ctx context.Context, course *domain.Course, userID int64) (bool, error) {

	if userID == course.CreatedBy {
		return true, nil
	}

	role, err := s.enrolledrepo.GetRoleUserInCourse(ctx, userID, course.Id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("%w: %v", errorsAPP.ErrInternalDB, err)
	}

	access, err := s.permissionrepo.HasPermission(ctx, role, domain.ResourceCourse, domain.ActionEdit)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}

		return false, fmt.Errorf("%w: %v", errorsAPP.ErrInternalDB, err)
	}

	return access, nil

}
