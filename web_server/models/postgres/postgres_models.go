package postgres_models

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresModels struct {
	pool *pgxpool.Pool
}

func (PGM PostgresModels) GetServerBySpaceId(space_id string) (int, error) {

	return 1, nil
}

func (PGM PostgresModels) GetServerByShareId(share_id string) (int, error) {

	return 1, nil
}
