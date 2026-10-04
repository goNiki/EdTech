package quiz

import (
	"context"
	"fmt"

	"edtech/internal/infrastructure/txmanager"
)

func (r *repositoryImpl) AcquireAdvisoryLock(ctx context.Context, userID int64, quizID int64) error {
	const op = "repository.quiz.AcquireAdvisoryLock"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)
	if _, err := q.Exec(ctx, "SELECT pg_advisory_xact_lock($1, $2)", userID, quizID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
