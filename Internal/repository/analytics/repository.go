package analytics

import (
	"edtech/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type repositoryAnalytics struct {
	db *pgxpool.Pool
}

func NewAnalyticsRepo(db *pgxpool.Pool) repository.AnalyticsRepository {
	return &repositoryAnalytics{
		db: db,
	}
}
