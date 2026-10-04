package txmanager

import (
	"context"
	"errors"
	"testing"

	"edtech/internal/infrastructure/db"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type mockTx struct {
	committed  bool
	rolledBack bool
	commitErr  error
	rbErr      error
}

func (m *mockTx) Begin(ctx context.Context) (pgx.Tx, error) {
	return m, nil
}

func (m *mockTx) Commit(ctx context.Context) error {
	m.committed = true
	return m.commitErr
}

func (m *mockTx) Rollback(ctx context.Context) error {
	m.rolledBack = true
	return m.rbErr
}

func (m *mockTx) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, nil
}

func (m *mockTx) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return nil, nil
}

func (m *mockTx) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return nil
}

func (m *mockTx) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	return 0, nil
}

func (m *mockTx) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	return nil
}

func (m *mockTx) LargeObjects() pgx.LargeObjects {
	return pgx.LargeObjects{}
}

func (m *mockTx) Prepare(ctx context.Context, name, sql string) (*pgconn.StatementDescription, error) {
	return nil, nil
}

func (m *mockTx) Conn() *pgx.Conn {
	return nil
}

type mockPool struct {
	beginCount int
	txToReturn *mockTx
	beginErr   error
}

func (m *mockPool) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	m.beginCount++
	if m.beginErr != nil {
		return nil, m.beginErr
	}
	return m.txToReturn, nil
}

type mockQueryExecutor struct {
	db.QueryExecutor
}

func TestTxManager_CommitSuccess(t *testing.T) {
	tx := &mockTx{}
	pool := &mockPool{txToReturn: tx}
	tm := NewTxManagerWithPool(pool)

	err := tm.WithTX(context.Background(), pgx.TxOptions{}, func(ctx context.Context) error {
		return nil
	})

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if pool.beginCount != 1 {
		t.Fatalf("expected BeginTx to be called once, got %d", pool.beginCount)
	}
	if !tx.committed {
		t.Fatalf("expected tx to be committed")
	}
	if tx.rolledBack {
		t.Fatalf("expected tx not to be rolled back")
	}
}

func TestTxManager_RollbackOnError(t *testing.T) {
	tx := &mockTx{}
	pool := &mockPool{txToReturn: tx}
	tm := NewTxManagerWithPool(pool)

	expectedErr := errors.New("custom business error")
	err := tm.WithTX(context.Background(), pgx.TxOptions{}, func(ctx context.Context) error {
		return expectedErr
	})

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected %v, got %v", expectedErr, err)
	}
	if !tx.rolledBack {
		t.Fatalf("expected tx to be rolled back")
	}
	if tx.committed {
		t.Fatalf("expected tx not to be committed")
	}
}

func TestTxManager_RollbackOnPanic(t *testing.T) {
	tx := &mockTx{}
	pool := &mockPool{txToReturn: tx}
	tm := NewTxManagerWithPool(pool)

	defer func() {
		p := recover()
		if p == nil {
			t.Fatalf("expected panic to be propagated")
		}
		if p != "something went terribly wrong" {
			t.Fatalf("unexpected panic value: %v", p)
		}
		if !tx.rolledBack {
			t.Fatalf("expected tx to be rolled back on panic")
		}
		if tx.committed {
			t.Fatalf("expected tx not to be committed")
		}
	}()

	_ = tm.WithTX(context.Background(), pgx.TxOptions{}, func(ctx context.Context) error {
		panic("something went terribly wrong")
	})
}

func TestTxManager_FlatNesting(t *testing.T) {
	tx := &mockTx{}
	pool := &mockPool{txToReturn: tx}
	tm := NewTxManagerWithPool(pool)

	nestedExecuted := false
	err := tm.WithTX(context.Background(), pgx.TxOptions{}, func(ctx context.Context) error {
		return tm.WithTX(ctx, pgx.TxOptions{}, func(innerCtx context.Context) error {
			nestedExecuted = true
			return nil
		})
	})

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !nestedExecuted {
		t.Fatalf("expected inner transaction to execute")
	}
	if pool.beginCount != 1 {
		t.Fatalf("expected BeginTx to be called exactly once for nested calls, got %d", pool.beginCount)
	}
	if !tx.committed {
		t.Fatalf("expected tx to be committed")
	}
}

func TestTxManager_GetQueryExecutor(t *testing.T) {
	tx := &mockTx{}
	pool := &mockPool{txToReturn: tx}
	tm := NewTxManagerWithPool(pool)
	defaultExecutor := &mockQueryExecutor{}

	qOut := GetQueryExecutor(context.Background(), defaultExecutor)
	if qOut != defaultExecutor {
		t.Fatalf("expected defaultExecutor when outside tx")
	}

	err := tm.WithTX(context.Background(), pgx.TxOptions{}, func(ctx context.Context) error {
		qIn := GetQueryExecutor(ctx, defaultExecutor)
		if qIn != tx {
			t.Fatalf("expected tx when inside tx context")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestWithTxContext_And_GetTxFromContext(t *testing.T) {
	ctx := context.Background()

	_, ok := GetTxFromContext(ctx)
	if ok {
		t.Fatalf("expected no tx in empty context")
	}

	mockT := &mockTx{}
	ctxWithTx := WithTxContext(ctx, mockT)

	tx, ok := GetTxFromContext(ctxWithTx)
	if !ok {
		t.Fatalf("expected tx in context")
	}
	if tx != mockT {
		t.Fatalf("expected returned tx to match mockTx")
	}
}

type limitedCapacityPool struct {
	maxActive  int
	activeChan chan struct{}
}

func newLimitedCapacityPool(maxActive int) *limitedCapacityPool {
	return &limitedCapacityPool{
		maxActive:  maxActive,
		activeChan: make(chan struct{}, maxActive),
	}
}

func (p *limitedCapacityPool) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	select {
	case p.activeChan <- struct{}{}:
		return &limitedTx{pool: p}, nil
	default:
		return nil, errors.New("connection pool exhausted: deadlock would occur")
	}
}

type limitedTx struct {
	mockTx
	pool     *limitedCapacityPool
	released bool
}

func (lt *limitedTx) Commit(ctx context.Context) error {
	if !lt.released {
		lt.released = true
		<-lt.pool.activeChan
	}
	return lt.mockTx.Commit(ctx)
}

func (lt *limitedTx) Rollback(ctx context.Context) error {
	if !lt.released {
		lt.released = true
		<-lt.pool.activeChan
	}
	return lt.mockTx.Rollback(ctx)
}

func TestTxManager_ReentrancyAndPoolStarvationPrevention(t *testing.T) {
	const poolCapacity = 5
	const concurrentWorkers = 20

	pool := newLimitedCapacityPool(poolCapacity)
	tm := NewTxManagerWithPool(pool)

	errChan := make(chan error, concurrentWorkers)

	for i := 0; i < concurrentWorkers; i++ {
		go func() {
			// Внешняя транзакция (например SubmitAttempt или GradeAttemptAnswer)
			err := tm.WithTX(context.Background(), pgx.TxOptions{}, func(ctx context.Context) error {
				// Вложенная транзакция (например CompleteLesson)
				return tm.WithTX(ctx, pgx.TxOptions{}, func(innerCtx context.Context) error {
					// Еще один уровень вложенности (например saveLessonProgress)
					return tm.WithTX(innerCtx, pgx.TxOptions{}, func(deepCtx context.Context) error {
						tx, ok := GetTxFromContext(deepCtx)
						if !ok || tx == nil {
							return errors.New("expected active tx in deep context")
						}
						return nil
					})
				})
			})
			errChan <- err
		}()
	}

	for i := 0; i < concurrentWorkers; i++ {
		if err := <-errChan; err != nil {
			t.Fatalf("worker failed with error (pool starvation or deadlock): %v", err)
		}
	}
}

