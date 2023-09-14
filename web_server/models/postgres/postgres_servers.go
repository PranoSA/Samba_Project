package postgres_models

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresServerModel struct {
	pool *pgxpool.Pool
}

func InitPostgresServerModel(pool *pgxpool.Pool) *PostgresServerModel {
	return &PostgresServerModel{pool}
}

func (PSM PostgresServerModel) GetServerBySpaceId(space_id string) (int, string, error) {

	sql := `
		SELECT Samba_File_Systems.server_id, Samba_Spaces.owner 
		FROM Samba_Spaces
		JOIN Samba_File_Systems
		ON Samba_File_Systems.fsid = Samba_Spaces.fs_id
		WHERE spaceid = @space_id
	`

	row, err := PSM.pool.Query(context.Background(), sql, pgx.NamedArgs{
		"space_id": space_id,
	})

	if err != nil {
		return 0, "", err
	}

	var serverid int
	var owner string
	row.Scan(&serverid)
	row.Scan(&owner)

	return serverid, owner, nil
}

func (PSM PostgresServerModel) GetServerByShareId(share_id string) (int, error) {

	sql := `
	SELECT space_id
	FROM Samba_Shares
	WHERE shareid = @shareid
	`

	row, err := PSM.pool.Query(context.Background(), sql, &pgx.NamedArgs{
		"shareid": share_id,
	})

	if err != nil {
		return 0, err
	}

	var space_id string
	row.Scan(&space_id)

	fsid, _, err := PSM.GetServerBySpaceId(space_id)
	if err != nil {
		return 0, err
	}

	return fsid, nil
}
