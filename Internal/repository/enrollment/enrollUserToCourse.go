package enrollment

import (
	"context"
	"edtech/internal/domain"
	"edtech/internal/infrastructure/db"
	errorsAPP "edtech/pkg/errors"
	"fmt"
)

func (r *repository) EnrollUserToCourse(ctx context.Context, q db.QueryExecutor, enroll domain.EnrolledInCourse) error {
	const op = "repository.enrolled.enrollusertocourse"

	query := `INSERT INTO users_courses (user_id, course_id, role) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`

	cmgTag, err := q.Exec(ctx, query, enroll.UserID, enroll.CourseID, enroll.Role)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if cmgTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, errorsAPP.ErrUserAlreadyEnrolled)
	}

	return nil
}
