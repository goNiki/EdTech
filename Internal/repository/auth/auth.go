package auth

import (
	"context"
	"edtech/internal/domain"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	Pool *pgxpool.Pool
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{
		Pool: pool,
	}
}
func (u *UserRepo) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	const op = "repository.user_repo.CreateUser"

	query := `INSERT INTO users (email, password_hash, role) VALUES ($1, $2, $3) RETURNING id , create_at`

	err := u.Pool.QueryRow(ctx, query, user.Email, user.PasswordHash, user.Role).Scan(&user.ID, &user.CreateAt)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil

}
