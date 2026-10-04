package enrollment

import (
	"context"
	"fmt"

	"edtech/internal/domain"
	"edtech/internal/infrastructure/txmanager"
	errorsAPP "edtech/pkg/errors"
)

func (r *repository) EnrollUserToCourse(ctx context.Context, enroll domain.EnrolledInCourse) error {
	const op = "repository.enrolled.enrollusertocourse"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	query := `
		INSERT INTO users_courses (user_id, course_id, role, enrolled_at) 
		VALUES ($1, $2, $3, NOW()) 
		ON CONFLICT (user_id, course_id) DO NOTHING
	`

	cmgTag, err := q.Exec(ctx, query, enroll.UserID, enroll.CourseID, enroll.Role)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if cmgTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, errorsAPP.ErrUserAlreadyEnrolled)
	}

	return nil
}
