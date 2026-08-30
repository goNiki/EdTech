package migrator

import (
	"database/sql"
	errorsAPP "edtech/pkg/errors"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

const tableName = "goose_migrations"

type Migrator struct {
	sqlDB  *sql.DB
	migDir string
}

func NewMigrator(pool *pgxpool.Pool, migDir string) (*Migrator, error) {
	sqlDB := stdlib.OpenDBFromPool(pool)

	if err := goose.SetDialect(string(goose.DialectPostgres)); err != nil {
		if cerr := sqlDB.Close(); cerr != nil {
			return nil, errors.Join(
				fmt.Errorf("%w: %w", errorsAPP.ErrSetGooseDialect, err),
				fmt.Errorf("%w: %w", errorsAPP.ErrCloseDb, cerr),
			)
		}
		return nil, fmt.Errorf("%w: %w", errorsAPP.ErrSetGooseDialect, err)
	}

	goose.SetTableName(tableName)

	return &Migrator{
		sqlDB:  sqlDB,
		migDir: migDir,
	}, nil
}

func (m *Migrator) Up() error {
	if err := goose.Up(m.sqlDB, m.migDir); err != nil {
		return fmt.Errorf("%w: %w", errorsAPP.ErrUpMigration, err)
	}
	return nil
}

func (m *Migrator) Down() error {
	if err := goose.Down(m.sqlDB, m.migDir); err != nil {
		return fmt.Errorf("%w: %w", errorsAPP.ErrDownMigration, err)
	}
	return nil
}

func (m *Migrator) DownTo(version int64) error {
	if err := goose.DownTo(m.sqlDB, m.migDir, version); err != nil {
		return fmt.Errorf("%w: %w", errorsAPP.ErrDownToMigration, err)
	}
	return nil
}

func (m *Migrator) Create(name string, migrationType string) error {
	if err := goose.Create(m.sqlDB, m.migDir, name, migrationType); err != nil {
		return fmt.Errorf("%w: %w", errorsAPP.ErrCreateMigration, err)
	}
	return nil
}

func (m *Migrator) Status() error {
	migration, err := goose.GetDBVersion(m.sqlDB)
	if err != nil {
		return fmt.Errorf("%w: %w", errorsAPP.ErrGetDbVersion, err)
	}
	log.Printf("Current database version: %d", migration)

	if err := goose.Status(m.sqlDB, m.migDir); err != nil {
		return fmt.Errorf("%w: %w", errorsAPP.ErrGetGooseStatus, err)
	}
	return nil
}

func (m *Migrator) UpTo(version int64) error {
	if err := goose.UpTo(m.sqlDB, m.migDir, version); err != nil {
		return fmt.Errorf("%w: %w", errorsAPP.ErrUpToMigration, err)
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
