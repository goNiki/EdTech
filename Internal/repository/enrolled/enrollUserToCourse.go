package enrolled

import (
	"context"
	"edtech/internal/domain"
	"fmt"
)

func (r *repository) EnrollUserToCourse(ctx context.Context, enroll domain.EnrolledInCourse) error {

	const op = "repository.enrolled.enrollusertocourse"

	query := `INSERT INTO users_courses (user_id, course_id, role) VALUES ($1, $2, $3)`

	_, err := r.Pool.Exec(ctx, query, enroll.UserID, enroll.CourseID, enroll.Role)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
