package analytics

import (
	"edtech/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type repositoryAnalytics struct {
	Pool *pgxpool.Pool
}

func NewAnalyticsRepo(pool *pgxpool.Pool) repository.AnalyticsRepository {
	return &repositoryAnalytics{
		Pool: pool,
	}
}
