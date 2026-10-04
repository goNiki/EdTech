package auth

import (
	"context"
	"errors"
	"fmt"

	"edtech/internal/domain"
	"edtech/internal/infrastructure/txmanager"
	repomodels "edtech/internal/repository/models"
	repoconverter "edtech/internal/repository/models/converter"
	errorsAPP "edtech/pkg/errors"

	"github.com/jackc/pgx/v5"
)

const baseUserSelect = `
	SELECT 
		id, 
		email, 
		password_hash, 
		username, 
		first_name, 
		last_name, 
		avatar_url, 
		bio, 
		headline, 
		role, 
		email_verified, 
		is_active, 
		is_banned, 
		last_login_at, 
		created_at, 
		updated_at, 
		deleted_at,
		preferences 
	FROM users
`

func scanUser(row pgx.Row) (*domain.User, error) {
	var user repomodels.User

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.AvatarURL,
		&user.Bio,
		&user.Headline,
		&user.Role,
		&user.EmailVerified,
		&user.IsActive,
		&user.IsBanned,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
		&user.Preferences,
	)
	if err != nil {
		return nil, err
	}

	return repoconverter.UserToDomain(&user), nil
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	const op = "repository.auth.GetUserByEmail"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	query := baseUserSelect + ` WHERE email = $1 AND deleted_at IS NULL`

	row := q.QueryRow(ctx, query, email)
	user, err := scanUser(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, errorsAPP.ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	return user, nil
}

func (r *repository) GetUserByUserName(ctx context.Context, username string) (*domain.User, error) {
	const op = "repository.auth.GetUserByUserName"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	query := baseUserSelect + ` WHERE username = $1 AND deleted_at IS NULL`

	row := q.QueryRow(ctx, query, username)
	user, err := scanUser(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, errorsAPP.ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	return user, nil
}

func (r *repository) GetUserByID(ctx context.Context, userID int64) (*domain.User, error) {
	const op = "repository.auth.GetUserByID"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	query := baseUserSelect + ` WHERE id = $1 AND deleted_at IS NULL`

	row := q.QueryRow(ctx, query, userID)
	user, err := scanUser(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, errorsAPP.ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	return user, nil
}
