package txmanager

import (
	"context"
	"errors"

	"edtech/internal/infrastructure/db"
	"github.com/jackc/pgx/v5"
)

type TransactionManager interface {
	WithTX(
		ctx context.Context, opts pgx.TxOptions,
		fn func(ctx context.Context, q db.QueryExecutor) error,
	) error
}

type TxManager struct {
	postgres *db.Postgres
}

func NewTxManager(postgres *db.Postgres) *TxManager {
	return &TxManager{
		postgres: postgres,
	}
}

func (tm *TxManager) WithTX(
	ctx context.Context, opts pgx.TxOptions,
	fn func(ctx context.Context, q db.QueryExecutor) error,
) error {
	tx, err := tm.postgres.Pool.BeginTx(ctx, opts)
	if err != nil {
		return err
	}

	fnErr := fn(ctx, tx)

	if fnErr != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil && !errors.Is(rbErr, pgx.ErrTxClosed) {
			return errors.Join(fnErr, rbErr)
		}
		return fnErr
	}

	return tx.Commit(ctx)
}
