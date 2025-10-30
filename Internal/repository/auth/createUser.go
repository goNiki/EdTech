package auth

import (
	"context"
	"edtech/internal/domain"
	"fmt"
)

func (r *repository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	const op = "repository.auth.CreateUser"

	query := `INSERT INTO users (email, password_hash, role, username) VALUES ($1, $2, $3, $4) RETURNING id , create_at`

	err := r.Pool.QueryRow(ctx, query, user.Email, user.PasswordHash, user.Role, user.Username).Scan(&user.ID, &user.CreateAt)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil

}
