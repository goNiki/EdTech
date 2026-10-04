package category

import (
	"edtech/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type repositoryImpl struct {
	Pool *pgxpool.Pool
}

func NewCategoryRepo(pool *pgxpool.Pool) repository.CategoryRepository {
	return &repositoryImpl{
		Pool: pool,
	}
}
