package refresh

import (
	"context"
	"edtech/internal/infrastructure/txmanager"
	"fmt"
	"time"
)

func (r *repository) Save(ctx context.Context, reftoken string, id int64, expiresAt time.Time) error {
	const op = "repository.refresh.Save"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	query := `INSERT INTO refresh_tokens (token, user_id, expires_at) VALUES ($1, $2, $3)`

	if _, err := q.Exec(ctx, query, reftoken, id, expiresAt); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
