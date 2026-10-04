package migrator

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	errorsAPP "edtech/pkg/errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

const tableName = "goose_migrations"

type Migrator struct {
	provider *goose.Provider
	sqlDB    *sql.DB
	migDir   string
}

func NewMigrator(pool *pgxpool.Pool, migDir string) (*Migrator, error) {
	sqlDB := stdlib.OpenDBFromPool(pool)

	provider, err := goose.NewProvider(
		goose.DialectPostgres,
		sqlDB,
		os.DirFS(migDir),
		goose.WithTableName(tableName),
	)
	if err != nil {
		if cerr := sqlDB.Close(); cerr != nil {
			return nil, errors.Join(
				fmt.Errorf("%w: %w", errorsAPP.ErrSetGooseDialect, err),
				fmt.Errorf("%w: %w", errorsAPP.ErrCloseDb, cerr),
			)
		}
		return nil, fmt.Errorf("%w: %w", errorsAPP.ErrSetGooseDialect, err)
	}

	return &Migrator{
		provider: provider,
		sqlDB:    sqlDB,
		migDir:   migDir,
	}, nil
}

func (m *Migrator) Up() error {
	if _, err := m.provider.Up(context.Background()); err != nil {
		return fmt.Errorf("%w: %w", errorsAPP.ErrUpMigration, err)
	}
	return nil
}

func (m *Migrator) Down() error {
	if _, err := m.provider.Down(context.Background()); err != nil {
		return fmt.Errorf("%w: %w", errorsAPP.ErrDownMigration, err)
	}
	return nil
}

func (m *Migrator) DownTo(version int64) error {
	if _, err := m.provider.DownTo(context.Background(), version); err != nil {
		return fmt.Errorf("%w: %w", errorsAPP.ErrDownToMigration, err)
	}
	return nil
}

func (m *Migrator) UpTo(version int64) error {
	if _, err := m.provider.UpTo(context.Background(), version); err != nil {
		return fmt.Errorf("%w: %w", errorsAPP.ErrUpToMigration, err)
	}
	return nil
}

func (m *Migrator) Status() error {
	v, err := m.provider.GetDBVersion(context.Background())
	if err != nil {
		return fmt.Errorf("%w: %w", errorsAPP.ErrGetDbVersion, err)
	}
	log.Printf("Current database version: %d", v)

	statuses, err := m.provider.Status(context.Background())
	if err != nil {
		return fmt.Errorf("%w: %w", errorsAPP.ErrGetGooseStatus, err)
	}
	for _, s := range statuses {
		log.Printf("Migration: %s, State: %s, AppliedAt: %v", s.Source.Path, s.State, s.AppliedAt)
	}
	return nil
}

func (m *Migrator) Create(name string, migrationType string) error {
	if err := goose.Create(m.sqlDB, m.migDir, name, migrationType); err != nil {
		return fmt.Errorf("%w: %w", errorsAPP.ErrCreateMigration, err)
	}
	return nil
}

func (m *Migrator) CloseDB() error {
	if m.sqlDB == nil {
		return nil
	}
	if err := m.sqlDB.Close(); err != nil {
		return fmt.Errorf("%w: %w", errorsAPP.ErrCloseDb, err)
	}
	return nil
}
