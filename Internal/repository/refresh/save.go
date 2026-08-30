package refresh

import (
	"context"
	"edtech/internal/infrastructure/db"
	"fmt"
	"time"
)

func (r *repository) Save(ctx context.Context, q db.QueryExecutor, reftoken string, id int64, expiresAt time.Time) error {
	const op = "repository.refresh.Save"

	query := `INSERT INTO refresh_tokens (token, user_id, expires_at) VALUES ($1, $2, $3)`

	if _, err := q.Exec(ctx, query, reftoken, id, expiresAt); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
