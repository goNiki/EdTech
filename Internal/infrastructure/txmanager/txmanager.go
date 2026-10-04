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

// WithTxContext помещает активную транзакцию pgx.Tx в context.Context.
func WithTxContext(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// GetTxFromContext извлекает активную транзакцию pgx.Tx из context.Context.
func GetTxFromContext(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(txKey{}).(pgx.Tx)
	return tx, ok
}

// WithTX открывает транзакцию или переиспользует существующую из контекста (реентерабельность).
// При наличии активной транзакции в ctx новая транзакция из пула НЕ открывается.
// При панике или ошибке выполняется безопасный Rollback.
func (tm *TxManager) WithTX(ctx context.Context, opts pgx.TxOptions, fn func(ctx context.Context) error) (err error) {
	if _, ok := GetTxFromContext(ctx); ok {
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

	ctxWithTx := WithTxContext(ctx, tx)
	if err = fn(ctxWithTx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func GetQueryExecutor(ctx context.Context, defaultPool db.QueryExecutor) db.QueryExecutor {
	if tx, ok := GetTxFromContext(ctx); ok {
		return tx
	}
	return defaultPool
}
