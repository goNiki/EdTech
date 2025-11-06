package enrolmend

import (
	"context"
	errorsAPP "edtech/pkg/errors"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (s *service) GetRoleUserInCource(ctx context.Context, userID, courseID int64) (string, error) {
	role, err := s.enrolledrepo.GetRoleUserInCourse(ctx, userID, courseID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil
		}
		return "", fmt.Errorf("%w: %v", errorsAPP.ErrInternalDB, err)
	}

	return role, nil
}
