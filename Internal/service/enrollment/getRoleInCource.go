package enrollment

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (s *service) GetRoleUserInCource(ctx context.Context, userID, courseID int64) (string, error) {
	const op = "service.enrollment.GetRoleUserInCource"

	role, err := s.enrolledrepo.GetRoleUserInCourse(ctx, s.db, userID, courseID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return role, nil
}
