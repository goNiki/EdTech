package review

import (
	"context"
	"errors"
	"fmt"

	"edtech/internal/domain"
	"edtech/internal/infrastructure/db"
	"edtech/internal/repository"
	errorsAPP "edtech/pkg/errors"

	"github.com/jackc/pgx/v5"
)

type repo struct{}

func NewReviewRepository() repository.ReviewRepository {
	return &repo{}
}

func (r *repo) UpsertReview(ctx context.Context, q db.QueryExecutor, review *domain.Review) (*domain.Review, error) {
	const op = "repository.review.UpsertReview"

	query := `
		INSERT INTO course_reviews (course_id, user_id, rating, comment, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		ON CONFLICT (course_id, user_id)
		DO UPDATE SET rating = EXCLUDED.rating, comment = EXCLUDED.comment, updated_at = NOW()
		RETURNING id, course_id, user_id, rating, comment, created_at, updated_at
	`

	var res domain.Review
	err := q.QueryRow(ctx, query, review.CourseID, review.UserID, review.Rating, review.Comment).
		Scan(&res.ID, &res.CourseID, &res.UserID, &res.Rating, &res.Comment, &res.CreatedAt, &res.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &res, nil
}

func (r *repo) DeleteReview(ctx context.Context, q db.QueryExecutor, courseID, userID int64) error {
	const op = "repository.review.DeleteReview"

	query := `DELETE FROM course_reviews WHERE course_id = $1 AND user_id = $2`

	tag, err := q.Exec(ctx, query, courseID, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, errorsAPP.ErrNotFound)
	}

	return nil
}

func (r *repo) GetReviewByUserAndCourse(ctx context.Context, q db.QueryExecutor, courseID, userID int64) (*domain.Review, error) {
	const op = "repository.review.GetReviewByUserAndCourse"

	query := `
		SELECT r.id, r.course_id, r.user_id, r.rating, r.comment, r.created_at, r.updated_at,
		       COALESCE(u.user_name, ''), COALESCE(u.avatar_url, '')
		FROM course_reviews r
		LEFT JOIN users u ON u.id = r.user_id
		WHERE r.course_id = $1 AND r.user_id = $2
	`

	var res domain.Review
	err := q.QueryRow(ctx, query, courseID, userID).
		Scan(&res.ID, &res.CourseID, &res.UserID, &res.Rating, &res.Comment, &res.CreatedAt, &res.UpdatedAt, &res.UserName, &res.UserAvatar)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorsAPP.ErrNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &res, nil
}

func (r *repo) ListReviewsByCourse(ctx context.Context, q db.QueryExecutor, courseID int64, limit, offset int) ([]domain.Review, int64, error) {
	const op = "repository.review.ListReviewsByCourse"

	countQuery := `SELECT COUNT(*) FROM course_reviews WHERE course_id = $1`
	var total int64
	if err := q.QueryRow(ctx, countQuery, courseID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("%s: count reviews: %w", op, err)
	}

	if total == 0 {
		return []domain.Review{}, 0, nil
	}

	query := `
		SELECT r.id, r.course_id, r.user_id, r.rating, r.comment, r.created_at, r.updated_at,
		       COALESCE(u.user_name, ''), COALESCE(u.avatar_url, '')
		FROM course_reviews r
		LEFT JOIN users u ON u.id = r.user_id
		WHERE r.course_id = $1
		ORDER BY r.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := q.Query(ctx, query, courseID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: select reviews: %w", op, err)
	}
	defer rows.Close()

	reviews := make([]domain.Review, 0, limit)
	for rows.Next() {
		var rev domain.Review
		if err := rows.Scan(
			&rev.ID,
			&rev.CourseID,
			&rev.UserID,
			&rev.Rating,
			&rev.Comment,
			&rev.CreatedAt,
			&rev.UpdatedAt,
			&rev.UserName,
			&rev.UserAvatar,
		); err != nil {
			return nil, 0, fmt.Errorf("%s: scan review: %w", op, err)
		}
		reviews = append(reviews, rev)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("%s: rows err: %w", op, err)
	}

	return reviews, total, nil
}

func (r *repo) GetCourseRatingSummary(ctx context.Context, q db.QueryExecutor, courseID int64) (float64, int, error) {
	const op = "repository.review.GetCourseRatingSummary"

	query := `
		SELECT COALESCE(AVG(rating), 0.0)::FLOAT8, COUNT(*)::INT
		FROM course_reviews
		WHERE course_id = $1
	`

	var avgRating float64
	var count int
	if err := q.QueryRow(ctx, query, courseID).Scan(&avgRating, &count); err != nil {
		return 0, 0, fmt.Errorf("%s: %w", op, err)
	}

	return avgRating, count, nil
}
