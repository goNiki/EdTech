package refresh

import (
	"context"
	"edtech/internal/infrastructure/db"
	errorsAPP "edtech/pkg/errors"
	"fmt"
)

func (r *repository) DeleteRefreshToken(ctx context.Context, q db.QueryExecutor, hashToken string) error {
	const op = "repository.refresh.DeleteRefreshToken"

	query := `
		DELETE FROM refresh_tokens
		WHERE token = $1
	`

	tag, err := q.Exec(ctx, query, hashToken)
	if err != nil {
		return fmt.Errorf("%s: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, errorsAPP.ErrRefreshTokenNotFound)
	}

	return nil
}

func (r *repository) DeleteAllByUserID(ctx context.Context, q db.QueryExecutor, userID int64) error {
	const op = "repository.refresh.DeleteAllByUserID"

	query := `
		DELETE FROM refresh_tokens
		WHERE user_id = $1
	`

	if _, err := q.Exec(ctx, query, userID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

