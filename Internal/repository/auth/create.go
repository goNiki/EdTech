package auth

import (
	"context"
	"edtech/internal/domain"
	"edtech/internal/infrastructure/db"
	repoconverter "edtech/internal/repository/models/converter"
	"fmt"
)

func (r *repository) CreateUser(ctx context.Context, q db.QueryExecutor, createUser domain.CreateUser) (*domain.User, error) {

	const op = "repository.auth.createuser"

	query := `
		INSERT INTO users (
			email, 
			password_hash, 
			username,
			role,  
			) VALUES ($1, $2, $3, $4) 
			RETURNING 
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
	`

	user := repoconverter.CreateUserToEntity(createUser)

	err := q.QueryRow(ctx, query,
		user.Email,
		user.PasswordHash,
		user.Username,
		user.Role,
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
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return repoconverter.UserToModel(&user), nil
}
