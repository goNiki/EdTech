package review

import (
	"context"
	"fmt"

	"edtech/internal/domain"
	"edtech/internal/infrastructure/txmanager"
	"edtech/internal/repository"
	services "edtech/internal/service"
	errorsAPP "edtech/pkg/errors"

	"github.com/jackc/pgx/v5"
)

type service struct {
	reviewRepo   repository.ReviewRepository
	courseRepo   repository.CourseRepository
	enrolledRepo repository.EnrolledRepository
	progressRepo repository.ProgressRepository
	txManager    txmanager.TransactionManager
}

func NewReviewService(
	reviewRepo repository.ReviewRepository,
	courseRepo repository.CourseRepository,
	enrolledRepo repository.EnrolledRepository,
	progressRepo repository.ProgressRepository,
	txManager txmanager.TransactionManager,
) services.ReviewServices {
	return &service{
		reviewRepo:   reviewRepo,
		courseRepo:   courseRepo,
		enrolledRepo: enrolledRepo,
		progressRepo: progressRepo,
		txManager:    txManager,
	}
}

func (s *service) AddOrUpdateReview(ctx context.Context, userID, courseID int64, rating int, comment *string) (*domain.Review, error) {
	const op = "service.review.AddOrUpdateReview"

	if rating < 1 || rating > 5 {
		return nil, fmt.Errorf("%s: %w: rating must be between 1 and 5", op, errorsAPP.ErrValidationFailed)
	}

	// 1. Verify that user is enrolled in the course
	isEnrolled, err := s.enrolledRepo.UserExistCourse(ctx, userID, courseID)
	if err != nil {
		return nil, fmt.Errorf("%s: check enrollment: %w", op, err)
	}
	if !isEnrolled {
		return nil, fmt.Errorf("%s: %w: user is not enrolled in this course", op, errorsAPP.ErrForbidden)
	}

	// 2. Verify that user has completed at least 30% of the course
	prog, err := s.progressRepo.GetCourseProgress(ctx, userID, courseID)
	if err != nil || prog == nil || prog.Percent < 30 {
		return nil, fmt.Errorf("%s: %w: minimum 30%% course progress required to leave a review", op, errorsAPP.ErrForbidden)
	}

	var savedReview *domain.Review

	// 3. Transactional Upsert and Course Rating Recalculation
	err = s.txManager.WithTX(ctx, pgx.TxOptions{}, func(ctx context.Context) error {
		rev, uErr := s.reviewRepo.UpsertReview(ctx, &domain.Review{
			CourseID: courseID,
			UserID:   userID,
			Rating:   rating,
			Comment:  comment,
		})
		if uErr != nil {
			return fmt.Errorf("upsert review: %w", uErr)
		}
		savedReview = rev

		avgRating, count, sErr := s.reviewRepo.GetCourseRatingSummary(ctx, courseID)
		if sErr != nil {
			return fmt.Errorf("calculate rating summary: %w", sErr)
		}

		if cErr := s.courseRepo.UpdateCourseRatingStats(ctx, courseID, avgRating, count); cErr != nil {
			return fmt.Errorf("update course rating: %w", cErr)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return savedReview, nil
}

func (s *service) DeleteReview(ctx context.Context, userID, courseID int64) error {
	const op = "service.review.DeleteReview"

	err := s.txManager.WithTX(ctx, pgx.TxOptions{}, func(ctx context.Context) error {
		if dErr := s.reviewRepo.DeleteReview(ctx, courseID, userID); dErr != nil {
			return fmt.Errorf("delete review: %w", dErr)
		}

		avgRating, count, sErr := s.reviewRepo.GetCourseRatingSummary(ctx, courseID)
		if sErr != nil {
			return fmt.Errorf("calculate rating summary: %w", sErr)
		}

		if cErr := s.courseRepo.UpdateCourseRatingStats(ctx, courseID, avgRating, count); cErr != nil {
			return fmt.Errorf("update course rating: %w", cErr)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *service) ListCourseReviews(ctx context.Context, courseID int64, page, pageSize int) (*domain.CourseReviewsSummary, error) {
	const op = "service.review.ListCourseReviews"

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	reviews, total, err := s.reviewRepo.ListReviewsByCourse(ctx, courseID, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	avgRating, count, err := s.reviewRepo.GetCourseRatingSummary(ctx, courseID)
	if err != nil {
		return nil, fmt.Errorf("%s: get rating summary: %w", op, err)
	}

	return &domain.CourseReviewsSummary{
		Reviews:       reviews,
		AverageRating: avgRating,
		ReviewsCount:  count,
		Total:         total,
	}, nil
}

func (s *service) GetMyReview(ctx context.Context, userID, courseID int64) (*domain.Review, error) {
	const op = "service.review.GetMyReview"

	review, err := s.reviewRepo.GetReviewByUserAndCourse(ctx, courseID, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return review, nil
}
