package quiz

import (
	"edtech/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repositoryImpl struct {
	Pool *pgxpool.Pool
}

func NewQuizRepo(pool *pgxpool.Pool) repository.QuizRepository {
	return &repositoryImpl{
		Pool: pool,
	}
}
