package db

import (
	"context"
	"edtech/internal/infrastructure/config"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	Pool *pgxpool.Pool
}

func New(cfg config.Postgres) (*Postgres, error) {
	// создаем контекст с таймаутом 10 секунд
	// defer cancel гарантирует что контекст будет отменен при выходе из функции
	// это необходимо чтобы подключение не "висело" вечно, если что-то пойдет не так
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// формируем data sourse name - строка для подключения к базе данных
	// Формат: postgres://user:password@host:port/database?sslmode=mode
	dsn := cfg.Address()
	// парсим DSN строку в конфиг пула

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parce config DB: %w", err)
	}

	// производим настройку пула

	poolConfig.MaxConns = cfg.MaxConns()                   // максимум соединений - ограничивает нагрузку на БД
	poolConfig.MinConns = cfg.MinConns()                   // минимум соединений в пуле - готовые соединения для быстрых запросов
	poolConfig.MaxConnLifetime = cfg.MaxConnLifeTime()     // Время жизни соединения - переподключение для свежих соединений
	poolConfig.MaxConnIdleTime = cfg.MaxConnIdleTime()     // время бездействия - освобождение неиспользуемых соединений
	poolConfig.HealthCheckPeriod = cfg.HealthCheckPeriod() // частота проверок - провряет работоспособность соединения

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Postgres{Pool: pool}, nil
}

func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
