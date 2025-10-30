package auth

import (
	"context"
	"edtech/internal/domain"
	errorsAPP "edtech/pkg/errors"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {

	const op = "repository.auth.GetUserByID"

	query := `SELECT id, email, password_hash, username, role, create_at FROM users WHERE email = $1`

	row := r.Pool.QueryRow(ctx, query, email)

	var user domain.User

	if err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Username, &user.Role, &user.CreateAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorsAPP.ErrUserNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}
