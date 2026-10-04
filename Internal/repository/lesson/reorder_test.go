package lesson_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"edtech/internal/infrastructure/txmanager"
	"edtech/internal/repository/lesson"
	errorsAPP "edtech/pkg/errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type recordedCall struct {
	SQL  string
	Args []any
}

type mockTx struct {
	execCalls []recordedCall
	execErr   error
}

func (m *mockTx) Begin(ctx context.Context) (pgx.Tx, error) { return m, nil }
func (m *mockTx) Commit(ctx context.Context) error          { return nil }
func (m *mockTx) Rollback(ctx context.Context) error        { return nil }
func (m *mockTx) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	m.execCalls = append(m.execCalls, recordedCall{SQL: sql, Args: arguments})
	if m.execErr != nil {
		return pgconn.CommandTag{}, m.execErr
	}
	return pgconn.NewCommandTag("UPDATE 1"), nil
}
func (m *mockTx) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) { return nil, nil }
func (m *mockTx) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row          { return nil }
func (m *mockTx) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	return 0, nil
}
func (m *mockTx) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults { return nil }
func (m *mockTx) LargeObjects() pgx.LargeObjects                               { return pgx.LargeObjects{} }
func (m *mockTx) Prepare(ctx context.Context, name, sql string) (*pgconn.StatementDescription, error) {
	return nil, nil
}
func (m *mockTx) Conn() *pgx.Conn { return nil }

type mockPool struct {
	tx *mockTx
}

func (m *mockPool) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	return m.tx, nil
}

func TestReorderLessons_Batch(t *testing.T) {
	repo := lesson.NewLessonRepo(nil)

	t.Run("empty slice does not execute queries", func(t *testing.T) {
		tx := &mockTx{}
		tm := txmanager.NewTxManagerWithPool(&mockPool{tx: tx})

		err := tm.WithTX(context.Background(), pgx.TxOptions{}, func(ctx context.Context) error {
			secID := int64(5)
			return repo.ReorderLessons(ctx, &secID, []int64{})
		})

		require.NoError(t, err)
		assert.Empty(t, tx.execCalls, "expected 0 SQL executions for empty lesson list")
	})

	t.Run("batch reorder with sectionID in single SQL query", func(t *testing.T) {
		tx := &mockTx{}
		tm := txmanager.NewTxManagerWithPool(&mockPool{tx: tx})

		secID := int64(42)
		lessonIDs := []int64{101, 102, 103, 104}

		err := tm.WithTX(context.Background(), pgx.TxOptions{}, func(ctx context.Context) error {
			return repo.ReorderLessons(ctx, &secID, lessonIDs)
		})

		require.NoError(t, err)
		require.Len(t, tx.execCalls, 1, "must execute exactly 1 batch query instead of N queries")

		call := tx.execCalls[0]
		assert.True(t, strings.Contains(call.SQL, "unnest($2::bigint[])"))
		assert.True(t, strings.Contains(call.SQL, "unnest($3::int[])"))
		require.Len(t, call.Args, 3)
		assert.Equal(t, int64(42), call.Args[0])
		assert.Equal(t, lessonIDs, call.Args[1])
		assert.Equal(t, []int32{1, 2, 3, 4}, call.Args[2])
	})

	t.Run("batch reorder without sectionID in single SQL query", func(t *testing.T) {
		tx := &mockTx{}
		tm := txmanager.NewTxManagerWithPool(&mockPool{tx: tx})

		lessonIDs := []int64{201, 202}

		err := tm.WithTX(context.Background(), pgx.TxOptions{}, func(ctx context.Context) error {
			return repo.ReorderLessons(ctx, nil, lessonIDs)
		})

		require.NoError(t, err)
		require.Len(t, tx.execCalls, 1, "must execute exactly 1 batch query")

		call := tx.execCalls[0]
		assert.True(t, strings.Contains(call.SQL, "unnest($1::bigint[])"))
		assert.True(t, strings.Contains(call.SQL, "unnest($2::int[])"))
		require.Len(t, call.Args, 2)
		assert.Equal(t, lessonIDs, call.Args[0])
		assert.Equal(t, []int32{1, 2}, call.Args[1])
	})

	t.Run("db error is wrapped with ErrInternalDB", func(t *testing.T) {
		tx := &mockTx{execErr: errors.New("db connection failure")}
		tm := txmanager.NewTxManagerWithPool(&mockPool{tx: tx})

		lessonIDs := []int64{301}

		err := tm.WithTX(context.Background(), pgx.TxOptions{}, func(ctx context.Context) error {
			return repo.ReorderLessons(ctx, nil, lessonIDs)
		})

		require.Error(t, err)
		assert.True(t, errors.Is(err, errorsAPP.ErrInternalDB))
	})
}
