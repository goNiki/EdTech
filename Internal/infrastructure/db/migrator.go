package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

type Migrator struct {
	pool   *pgxpool.Pool
	migDir string
}

func NewMigrator(pool *pgxpool.Pool, migrationsDir string) *Migrator {
	return &Migrator{
		pool:   pool,
		migDir: migrationsDir,
	}
}

func (m *Migrator) getDB() (*sql.DB, error) {
	conn, err := m.pool.Acquire(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to acquire connection, %w", err)
	}

	defer conn.Release()

	sqlDB := stdlib.OpenDBFromPool(m.pool)
	return sqlDB, nil
}

func (m *Migrator) Up() error {
	sqlDB, err := m.getDB()
	if err != nil {
		return err
	}
	defer func() {
		if cerr := closeDB(sqlDB); err != nil {
			err = errors.Join(err, cerr)
		}
	}()

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	goose.SetTableName("goose_migrations")

	if err := goose.Up(sqlDB, m.migDir); err != nil {
		return fmt.Errorf("failet to up migrations: %w", err)
	}
	return nil
}

func (m *Migrator) Down() error {

	sqlDB, err := m.getDB()
	if err != nil {
		return err
	}
	defer func() {
		if cerr := closeDB(sqlDB); err != nil {
			err = errors.Join(err, cerr)
		}
	}()

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	if err = goose.Down(sqlDB, m.migDir); err != nil {
		return fmt.Errorf("failed to down migration: %w", err)
	}
	return nil

}

func (m *Migrator) DownTo(version int64) error {
	sqlDB, err := m.getDB()
	if err != nil {
		return err
	}
	defer func() {
		if cerr := closeDB(sqlDB); err != nil {
			err = errors.Join(err, cerr)
		}
	}()

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	if err = goose.DownTo(sqlDB, m.migDir, version); err != nil {
		return fmt.Errorf("failed to downto migration: %w", err)
	}
	return nil
}

func (m *Migrator) Create(name string, migrationType string) error {
	sqlDB, err := m.getDB()
	if err != nil {
		return nil
	}

	defer func() {
		if cerr := closeDB(sqlDB); err != nil {
			err = errors.Join(err, cerr)
		}
	}()

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	if err = goose.Create(sqlDB, m.migDir, name, migrationType); err != nil {
		return fmt.Errorf("failed to create migration: %w", err)
	}

	return nil
}

func (m *Migrator) Status() error {

	sqlDB, err := m.getDB()

	if err != nil {
		return err
	}

	defer func() {
		if cerr := closeDB(sqlDB); err != nil {
			err = errors.Join(err, cerr)
		}
	}()

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	migration, err := goose.GetDBVersion(sqlDB)
	if err != nil {
		return fmt.Errorf("failed GetDB Version: %w", err)
	}
	log.Printf("Current database version: %d", migration)

	if err = goose.Status(sqlDB, m.migDir); err != nil {
		return fmt.Errorf("faild to get status mifretions: %w", err)
	}

	return nil
}

func (m *Migrator) UpTo(version int64) error {
	sqlDB, err := m.getDB()

	if err != nil {
		return err
	}

	defer func() {
		if cerr := closeDB(sqlDB); err != nil {
			err = errors.Join(err, cerr)
		}
	}()

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	if err = goose.UpTo(sqlDB, m.migDir, version); err != nil {
		return fmt.Errorf("failed to UpTo migration: %w", err)
	}

	return nil
}

func closeDB(db *sql.DB) error {
	if db == nil {
		return nil
	}
	if err := db.Close(); err != nil {
		return fmt.Errorf("failed to close DB connection: %w", err)
	}
	return nil
}
