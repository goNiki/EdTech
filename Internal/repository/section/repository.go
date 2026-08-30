package section

import (
	"edtech/internal/infrastructure/db"
	"edtech/internal/repository"
)

type repositorySection struct {
	db db.QueryExecutor
}

func NewSectionRepository(database db.QueryExecutor) repository.SectionRepository {
	return &repositorySection{
		db: database,
	}
}
