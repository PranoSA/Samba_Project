package postgres_models

import (
	"context"
	"errors"
	"fmt"
	"time"

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

func (PGM PostgresSpaceModel) GetServerBySpaceId(space_id string) (int, string, error) {

	sql := `
		SELECT Samba_File_Systems.server_id, Samba_Spaces.owner 
		FROM Samba_Spaces
		JOIN Samba_File_Systems
		ON Samba_File_Systems.fsid = Samba_Spaces.fs_id
		WHERE spaceid = @space_id
	`

	row, err := PGM.pool.Query(context.Background(), sql, pgx.NamedArgs{
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := c.AlloateSpace(ctx, &proto_samba_management.SpaceAllocationRequest{
		Owner: ssr.Owner,
		Size:  ssr.Megabytes,
	})

	if err == context.DeadlineExceeded {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

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

	serverid, _, e := PSM.GetServerBySpaceId(dsr.Space_id)
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

	sql := `
	SELECT owner, spaceid, owner, alloc_size
	FROM Samba_Spaces
	WHERE owner = @owner
	`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	rows, err := PSM.pool.Query(ctx, sql, pgx.NamedArgs{
		"owner": owner_uuid,
	})

	defer rows.Close()

	if err != nil {
		return nil, err
	}

	var response []models.SpaceResponse
	for rows.Next() {
		var NextSpaceResponse models.SpaceResponse = models.SpaceResponse{}

		scanerr := rows.Scan(&NextSpaceResponse.Email, &NextSpaceResponse.Spaceid, &NextSpaceResponse.Owner, &NextSpaceResponse.Megabytes)
		if scanerr != nil {
			return nil, err
		}
		response = append(response, NextSpaceResponse)
	}

	return &response, nil //&[]models.SpaceResponse{}, nil
}
