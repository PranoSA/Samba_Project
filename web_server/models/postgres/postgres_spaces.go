package postgres_models

import (
	"context"
	"errors"
	"fmt"

	"github.com/PranoSA/samba_share_backend/proto_samba_management"
	"github.com/PranoSA/samba_share_backend/web_server/grpc_webclient"
	"github.com/PranoSA/samba_share_backend/web_server/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresSpaceModel struct {
	pool *pgxpool.Pool
}

func InitPostgresSpaceModel(pool *pgxpool.Pool) *PostgresSpaceModel {
	return &PostgresSpaceModel{pool: pool}
}

func (PGM PostgresSpaceModel) GetServerBySpaceId(space_id string) (int, error) {

	sql := `
		SELECT server_id FROM Samba_Spaces
		WHERE spaceid = @space_id 
		
	`

	row, err := PGM.pool.Query(context.Background(), sql, pgx.NamedArgs{
		"space_id": space_id,
	})

	if err != nil {
		return 0, err
	}

	var serverid int

	row.Scan(&serverid)

	return serverid, nil
}

func (PGM PostgresSpaceModel) GetServerByShareId(share_id string) (int, error) {

	return 1, nil
}

func (PSM PostgresSpaceModel) CreateSpace(ssr models.SpaceRequest) (*models.SpaceResponse, error) {

	server_id := grpc_webclient.GetAndUpdateNextId()

	/**
	 *
	 * Process SQL QUERY
	 *
	 */

	c := grpc_webclient.GRPCSambaClients[server_id].Grpc_Samba_Client

	res, _ := c.AlloateSpace(context.Background(), &proto_samba_management.SpaceAllocationRequest{
		Owner: ssr.Owner,
		Size:  ssr.Megabytes,
	})

	sql := `
	INSERT INTO Samba_Spaces (fs_id, owner, size)
	VALUES(@fs_id, @owner, @size)
`

	row, err := PSM.pool.Query(context.Background(), sql, pgx.NamedArgs{
		"fs_id": res.Fsid,
		"owner": ssr.Owner,
		"size":  res.Size,
	})

	if err != nil {

	}

	row.Scan(nil)

	return &models.SpaceResponse{}, nil
}

func (PSM PostgresSpaceModel) DeleteSpaceById(dsr models.DeleteSpaceRequest) (*models.SpaceResponse, error) {

	serverid, e := PSM.GetServerBySpaceId(dsr.Space_id)
	if e != nil {
		return nil, e
	}
	//Delete Space Here

	sql := `
		DELETE FROM Samba_Spaces
		WHERE spaceid = @spaceid
		AND owner = @owner
	`

	row, err := PSM.pool.Query(context.Background(), sql, &pgx.NamedArgs{
		"spaceid": dsr.Space_id,
		"owner":   dsr.Owner,
	})

	if err != nil {
		return nil, errors.New("Entry Doesn't Exist In Database")
	}

	defer row.Close()

	c := grpc_webclient.GRPCSambaClients[serverid].Grpc_Samba_Client
	c.DeleteSpace(context.Background(), &proto_samba_management.DeleteSpaceRequest{
		Spaceid: dsr.Space_id,
	})

	fmt.Println(c)

	return &models.SpaceResponse{}, nil
}

func (PSM PostgresSpaceModel) GetSpaceById(models.DeleteSpaceRequest) (*models.SpaceResponse, error) {

	return &models.SpaceResponse{}, nil
}

func (PSM PostgresSpaceModel) GetSpaceByOwner(owner_uuid string) (*[]models.SpaceResponse, error) {

	return &[]models.SpaceResponse{}, nil
}
