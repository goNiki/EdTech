package auth

import (
	"context"
	"edtech/internal/infrastructure/txmanager"
	"fmt"
)

func (r *repository) ExistingByUsernName(ctx context.Context, username string) (bool, error) {
	const op = "repository.auth.existingbyusername"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	query := `
		SELECT EXISTS(
			SELECT 1 
			FROM user
			WHERE username = $1 AND deleted_at IS NULL
		)
	`
	var exists bool

	err := q.QueryRow(ctx, query, username).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return exists, nil

}

func (r *repository) ExistingByEmail(ctx context.Context, email string) (bool, error) {
	const op = "repository.auth.existingbyemail"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	query := `
		SELECT EXISTS(
			SELECT 1 
			FROM user
			WHERE email = $1 AND deleted_at IS NULL
		)
	`
	var exists bool

	err := q.QueryRow(ctx, query, email).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return exists, nil

}
