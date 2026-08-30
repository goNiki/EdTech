package auth

import (
	"context"
	"edtech/internal/domain"
	"edtech/internal/infrastructure/db"
	errorsAPP "edtech/pkg/errors"
	"fmt"
	"time"
)

func (r *repository) UpdateLastLogin(ctx context.Context, q db.QueryExecutor, userID int64, now time.Time) error {
	const op = "repository.auth.UpdateLastLogin"

	query := `
		UPDATE users
		SET last_login_at = $2
		WHERE id = $1 AND deleted_at IS NULL
	`
	if _, err := q.Exec(
		ctx,
		query,
		userID,
		now,
	); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *repository) UpdateProfile(ctx context.Context, q db.QueryExecutor, user *domain.User) error {
	const op = "repository.auth.UpdateProfile"

	query := `
		UPDATE users
		SET first_name = $2,
		    last_name = $3,
		    avatar_url = $4,
		    bio = $5,
		    updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	tag, err := q.Exec(
		ctx,
		query,
		user.ID,
		user.FirstName,
		user.LastName,
		user.AvatarURL,
		user.Bio,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if tag.RowsAffected() == 0 {
		return errorsAPP.ErrUserNotFound
	}

	return nil
}

func (r *repository) UpdatePassword(ctx context.Context, q db.QueryExecutor, userID int64, passHash string) error {
	const op = "repository.auth.UpdatePassword"

	query := `
		UPDATE users
		SET password_hash = $2,
			updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	tag, err := q.Exec(ctx, query, userID, passHash)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if tag.RowsAffected() == 0 {
		return errorsAPP.ErrUserNotFound
	}

	return nil
}

func (r *repository) SetEmailVerified(ctx context.Context, q db.QueryExecutor, userID int64, verified bool) error {
	const op = "repository.auth.SetEmailVerified"

	query := `
		UPDATE users
		SET email_verified = $2,
			updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	tag, err := q.Exec(ctx, query, userID, verified)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if tag.RowsAffected() == 0 {
		return errorsAPP.ErrUserNotFound
	}

	return nil
}

func (r *repository) UpdateRole(ctx context.Context, q db.QueryExecutor, userID int64, role domain.Role) error {
	const op = "repository.auth.UpdateRole"

	query := `
		UPDATE users
		SET role = $2,
			updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	tag, err := q.Exec(ctx, query, userID, string(role))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if tag.RowsAffected() == 0 {
		return errorsAPP.ErrUserNotFound
	}

	return nil
}

func (r *repository) SetBannedStatus(ctx context.Context, q db.QueryExecutor, userID int64, isBanned bool) error {
	const op = "repository.auth.SetBannedStatus"

	query := `
		UPDATE users
		SET is_banned = $2,
			updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	tag, err := q.Exec(ctx, query, userID, isBanned)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if tag.RowsAffected() == 0 {
		return errorsAPP.ErrUserNotFound
	}

	return nil
}
