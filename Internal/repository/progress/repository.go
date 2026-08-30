package progress

import (
	repo "edtech/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

var _ repo.ProgressRepository = (*repository)(nil)

type repository struct {
	Pool *pgxpool.Pool
}

func NewProgressRepo(pool *pgxpool.Pool) *repository {
	return &repository{
		Pool: pool,
	}
}
