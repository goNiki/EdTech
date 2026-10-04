package section

import (
	"edtech/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type repositorySection struct {
	Pool *pgxpool.Pool
}

func NewSectionRepository(pool *pgxpool.Pool) repository.SectionRepository {
	return &repositorySection{
		Pool: pool,
	}
}
