package txmanager

import (
	"context"

	"edtech/internal/infrastructure/db"
	"github.com/jackc/pgx/v5"
)

type txKey struct{}

type TransactionManager interface {
	WithTX(ctx context.Context, opts pgx.TxOptions, fn func(ctx context.Context) error) error
}

type PgxPool interface {
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

type TxManager struct {
	pool PgxPool
}

func NewTxManager(postgres *db.Postgres) *TxManager {
	return &TxManager{
		pool: postgres.Pool,
	}
}

func NewTxManagerWithPool(pool PgxPool) *TxManager {
	return &TxManager{
		pool: pool,
	}
}

func (tm *TxManager) WithTX(ctx context.Context, opts pgx.TxOptions, fn func(ctx context.Context) error) (err error) {
	if _, ok := ctx.Value(txKey{}).(pgx.Tx); ok {
		return fn(ctx)
	}

	tx, err := tm.pool.BeginTx(ctx, opts)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(context.Background())
			panic(p)
		} else if err != nil {
			_ = tx.Rollback(context.Background())
		}
	}()

	ctxWithTx := context.WithValue(ctx, txKey{}, tx)
	if err = fn(ctxWithTx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func GetQueryExecutor(ctx context.Context, defaultPool db.QueryExecutor) db.QueryExecutor {
	if tx, ok := ctx.Value(txKey{}).(pgx.Tx); ok {
		return tx
	}
	return defaultPool
}
