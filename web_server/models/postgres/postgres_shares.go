package postgres_models

import (
	"context"

	"github.com/PranoSA/samba_share_backend/proto_samba_management"
	"github.com/PranoSA/samba_share_backend/web_server/grpc_webclient"
	"github.com/PranoSA/samba_share_backend/web_server/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresShareModel struct {
	pool *pgxpool.Pool
}

func InitPostgresShareModel(pool *pgxpool.Pool) *PostgresShareModel {
	return &PostgresShareModel{pool: pool}
}

/**
 * type SambaShareModel interface {
	AddShare(SambaShareResponse) (SambaShareResponse, error)
	DeleteShare(SambaShareResponse) (SambaShareResponse, error)
	CreateInvite(ShareInviteRequest) (ShareInviteResponse, error)
	AcceptInvite(ShareInviteAccept) (ShareInviteResponse, error)
}

*/

func (PGM PostgresShareModel) GetServerBySpaceId(space_id string) (int, error) {

	return 1, nil
}

func (PGM PostgresShareModel) GetServerByShareId(share_id string) (int, error) {

	return 1, nil
}

func (PShareM PostgresShareModel) AddShare(ssr models.SambaShareRequest) (*models.SambaShareResponse, error) {

	server_id, err := PShareM.GetServerBySpaceId(ssr.Spaceid)
	if err != nil {
		return nil, err
	}

	sql := `
		INSERT INTO Samba_Shares ()
		VALUES ()
		RETURNING shareid 
	`

	rows, err := PShareM.pool.Query(context.Background(), sql, pgx.NamedArgs{
		"space_id": ssr.Spaceid,
	})

	var shareid string

	rows.Scan(&shareid)

	if err != nil {
		return nil, err
	}

	res, err := grpc_webclient.GRPCSambaClients[server_id].Grpc_Samba_Client.AllocateSambaShare(context.Background(), &proto_samba_management.RequestShambaShare{
		Owner:   ssr.Email,
		Spaceid: ssr.Spaceid,
	})

	if err != nil {
		return nil, err
	}

	return &models.SambaShareResponse{
		Email:   ssr.Email,
		Shareid: res.Fsid,
	}, nil

	//If Succeeds

}
