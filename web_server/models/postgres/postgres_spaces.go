package postgres_models

import (
	"context"

	"github.com/PranoSA/samba_share_backend/proto_samba_management"
	"github.com/PranoSA/samba_share_backend/web_server/grpc_webclient"
	"github.com/PranoSA/samba_share_backend/web_server/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresSpaceModel struct {
	pool *pgxpool.Pool
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
	c := grpc_webclient.GRPCSambaClients[server_id].GRPC_Space_Client

	c.AlloateSpace(context.Background(), &proto_samba_management.SpaceAllocationRequest{
		Owner: ssr.Owner,
		Size:  ssr.Megabytes,
	})

	return &models.SpaceResponse{}, nil
}

func (PSM PostgresSpaceModel) DeleteSpaceById(id string) (*models.SpaceResponse, error) {

	_, e := PSM.GetServerBySpaceId(id)
	if e != nil {
		return nil, e
	}
	//Delete Space Here
	//c := grpc_webclient.GRPCSambaClients[server_id].GRPC_Space_Client

}

func (PSM PostgresSpaceModel) GetSpaceById(string) (*models.SpaceResponse, error) {

	return &models.SpaceResponse{}, nil
}
