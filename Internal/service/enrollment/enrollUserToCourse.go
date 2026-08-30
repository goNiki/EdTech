package enrollment

import (
	"context"
	"fmt"

	"edtech/internal/domain"
	"edtech/internal/infrastructure/db"
	errorsAPP "edtech/pkg/errors"

	"github.com/jackc/pgx/v5"
)

func (s *service) SelfEnrollCourse(ctx context.Context, req domain.SelfEnrollRequest) error {
	const op = "service.enrollment.SelfEnrollCourse"

	course, err := s.courserepo.GetCourseByID(ctx, s.db, req.CourseID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := course.CanSelfEnroll(); err != nil {
		return fmt.Errorf("%s: %w", op, errorsAPP.ErrForbidden)
	}

	enroll := domain.EnrolledInCourse{
		UserID:   req.UserID,
		CourseID: req.CourseID,
		Role:     string(domain.RoleStudent),
	}

	err = s.executeEnrollmentTransaction(ctx, enroll)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *service) TeacherEnrollCourse(ctx context.Context, req domain.TeacherEnrollRequest) error {
	const op = "service.enrollment.TeacherEnrollCourse"

	course, err := s.courserepo.GetCourseByID(ctx, s.db, req.CourseID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	canManage, err := s.accessService.CanManageCourseUsers(ctx, course, req.TeacherID)
	if err != nil || !canManage {
		return fmt.Errorf("%s: %w", op, errorsAPP.ErrForbidden)
	}

	if req.Role == string(domain.RoleStudent) && course.Status != domain.StatusPublished {
		return fmt.Errorf("%s: %w", op, errorsAPP.ErrCannotEnrollStudentInDraft)
	}

	targetUser, err := s.userrepo.GetUserByEmail(ctx, s.db, req.TargetEmail)
	if err != nil {
		return fmt.Errorf("%s: target user not found: %w", op, err)
	}

	enroll := domain.EnrolledInCourse{
		UserID:   targetUser.ID,
		CourseID: req.CourseID,
		Role:     req.Role,
	}

	return s.executeEnrollmentTransaction(ctx, enroll)
}

func (s *service) executeEnrollmentTransaction(ctx context.Context, enroll domain.EnrolledInCourse) error {
	err := s.txManager.WithTX(ctx, pgx.TxOptions{}, func(ctx context.Context, tx db.QueryExecutor) error {
		if err := s.enrolledrepo.EnrollUserToCourse(ctx, tx, enroll); err != nil {
			return err
		}
		if err := s.courserepo.IncrementEnrolledCount(ctx, tx, enroll.CourseID); err != nil {
			return err
		}
		return nil
	})

	return err
}
