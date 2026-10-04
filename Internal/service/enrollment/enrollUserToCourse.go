package enrollment

import (
	"context"
	"fmt"

	"edtech/internal/domain"
	errorsAPP "edtech/pkg/errors"

	"github.com/jackc/pgx/v5"
)

func (s *service) SelfEnrollCourse(ctx context.Context, req domain.SelfEnrollRequest) error {
	const op = "service.enrollment.SelfEnrollCourse"

	course, err := s.courserepo.GetCourseByID(ctx, req.CourseID)
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

	if err := s.executeEnrollmentTransaction(ctx, enroll); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *service) TeacherEnrollCourse(ctx context.Context, req domain.TeacherEnrollRequest) error {
	const op = "service.enrollment.TeacherEnrollCourse"

	course, err := s.courserepo.GetCourseByID(ctx, req.CourseID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	canManage, err := s.accessService.CanManageCourseUsers(ctx, course, req.TeacherID)
	if err != nil || !canManage {
		return fmt.Errorf("%s: %w", op, errorsAPP.ErrForbidden)
	}

	var targetUserID int64
	if req.TargetUserID != nil && *req.TargetUserID > 0 {
		targetUserID = *req.TargetUserID
		_, err := s.userrepo.GetUserByID(ctx, targetUserID)
		if err != nil {
			return fmt.Errorf("%s: target user not found: %w", op, err)
		}
	} else if req.TargetEmail != "" {
		targetUser, err := s.userrepo.GetUserByEmail(ctx, req.TargetEmail)
		if err != nil {
			return fmt.Errorf("%s: target user with email '%s' not found: %w", op, req.TargetEmail, err)
		}
		targetUserID = targetUser.ID
	} else {
		return fmt.Errorf("%s: %w: email or user_id is required", op, errorsAPP.ErrValidationFailed)
	}

	role := req.Role
	if role == "" {
		role = string(domain.RoleStudent)
	}

	enroll := domain.EnrolledInCourse{
		UserID:   targetUserID,
		CourseID: req.CourseID,
		Role:     role,
	}

	if err := s.executeEnrollmentTransaction(ctx, enroll); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *service) executeEnrollmentTransaction(ctx context.Context, enroll domain.EnrolledInCourse) error {
	return s.txManager.WithTX(ctx, pgx.TxOptions{}, func(ctx context.Context) error {
		if err := s.enrolledrepo.EnrollUserToCourse(ctx, enroll); err != nil {
			return err
		}
		return s.courserepo.IncrementEnrolledCount(ctx, enroll.CourseID)
	})
}
