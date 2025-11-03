package permission

import "github.com/jackc/pgx/v5/pgxpool"

type repository struct {
	Pool *pgxpool.Pool
}

func NewPermissionRepo(pool *pgxpool.Pool) *repository {
	return &repository{
		Pool: pool,
	}
}
