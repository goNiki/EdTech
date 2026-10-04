package auth

import (
	"context"
	"edtech/internal/domain"
	"edtech/internal/infrastructure/txmanager"
	repomodels "edtech/internal/repository/models"
	repoconverter "edtech/internal/repository/models/converter"
	errorsAPP "edtech/pkg/errors"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	const op = "repository.auth.getuserbyemail"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	query := `SELECT id, email, password_hash, username, first_name, last_name, avatar_url, bio, role, email_verified, is_active, is_banned, last_login_at, created_at, updated_at, deleted_at FROM users WHERE email = $1 AND deleted_at IS NULL`

	var user repomodels.User

	err := q.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.AvatarURL,
		&user.Bio,
		&user.Role,
		&user.EmailVerified,
		&user.IsActive,
		&user.IsBanned,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, errorsAPP.ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	return repoconverter.UserToDomain(&user), nil
}

func (r *repository) GetUserByUserName(ctx context.Context, username string) (*domain.User, error) {
	const op = "repository.auth.getuserbyusername"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	query := `
		SELECT 
			id, 
			email, 
			password_hash, 
			username, 
			first_name, 
			last_name, 
			avatar_url, 
			bio, 
			role, 
			email_verified, 
			is_active, 
			is_banned, 
			last_login_at, 
			created_at, 
			updated_at, 
			deleted_at 
		FROM users 
		WHERE username = $1 AND deleted_at IS NULL
	`

	var user repomodels.User

	err := q.QueryRow(
		ctx,
		query,
		username,
	).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.AvatarURL,
		&user.Bio,
		&user.Role,
		&user.EmailVerified,
		&user.IsActive,
		&user.IsBanned,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, errorsAPP.ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	return repoconverter.UserToDomain(&user), nil

}

func (r *repository) GetUserByID(ctx context.Context, userID int64) (*domain.User, error) {
	const op = "repository.auth.getuserbyid"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	query := `
		SELECT 
			id, 
			email, 
			password_hash, 
			username, 
			first_name, 
			last_name, 
			avatar_url, 
			bio, 
			role, 
			email_verified, 
			is_active, 
			is_banned, 
			last_login_at, 
			created_at, 
			updated_at, 
			deleted_at 
		FROM users 
		WHERE id = $1 AND deleted_at IS NULL
	`

	var user repomodels.User

	err := q.QueryRow(
		ctx,
		query,
		userID,
	).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.AvatarURL,
		&user.Bio,
		&user.Role,
		&user.EmailVerified,
		&user.IsActive,
		&user.IsBanned,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, errorsAPP.ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	return repoconverter.UserToDomain(&user), nil

}
