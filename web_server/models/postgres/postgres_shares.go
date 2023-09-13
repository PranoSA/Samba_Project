package postgres_models

import (
	"context"
	"errors"
	"time"

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

func (PGM PostgresShareModel) DeleteShare(ssr models.SambaShareResponse) (*models.SambaShareResponse, error) {
	return nil, nil
}

func (PGM PostgresShareModel) AcceptInvite(sia models.ShareInviteAccept) (*models.ShareInviteResponse, error) {

	return nil, nil
}

func (PGM PostgresShareModel) CreateInvite(sir models.ShareInviteRequest) (*models.ShareInviteResponse, error) {

	invite, hash, expir := models.GenInvite()

	sql := `
	INSERT INTO Samba_Invites(share_id, owner, time_created, time_expired, invite_code, hash_code)
	VALUES(@share_id,@owner,@time_created,@time_expired,@invite_code,@hash_code)
	JOIN Samba_Shares
	WHERE Samba_Shares.Owner = @share_Owner
	RETURNING inviteid
	`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer cancel()

	txid, err := PGM.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})

	if err != nil {
	}

	sqlEnsureOwner := `
	SELECT shareid
	FROM Samba_Shares
	WHERE owner = @owner AND shareid = @shareid
	`

	_, err = txid.Query(ctx, sqlEnsureOwner, pgx.NamedArgs{
		"owner":   sir.Email,
		"shareid": sir.Shareid,
	})

	if err == pgx.ErrNoRows {
		return nil, models.ErrorEntryDoesNotExist
	}

	if err != nil {
		return nil, models.ErrorDatabaseTImeout
	}

	row, err := txid.Query(ctx, sql, pgx.NamedArgs{
		"share_id":     sir.Shareid,
		"owner":        sir.Email,
		"time_created": time.Now(),
		"time_expired": expir,
		"invite_code":  invite,
		"hash_code":    hash,
		"share_owner":  sir.Email,
	})

	if err == context.DeadlineExceeded {
		txid.Rollback(ctx)
		return nil, models.ErrorDatabaseTImeout
	}

	if err != nil {
		txid.Rollback(ctx)
		return nil, models.ErrorEntryDoesNotExist
	}

	var inviteid string
	row.Scan(&inviteid)

	txid.Commit(ctx)

	return &models.ShareInviteResponse{
		Email:       sir.Email,
		Inviteid:    inviteid,
		Invite_code: string(hash), // This will be utter rubbish
	}, nil

}

func (PGM PostgresShareModel) GetServerBySpaceId(space_id string) (int, string, error) {

	sql := `
	SELECT Samba_File_Systems.server_id, Samba_Spaces.owner 
	FROM Samba_Spaces
	JOIN Samba_File_Systems
	ON Samba_File_Systems.fsid = Samba_Spaces.fs_id
	WHERE spaceid = @space_id
	`
	//'6ff8278c-f1c5-458d-8d72-145851e5712a'
	//`
	//+ fmt.Sprintf("'lalalalal%s'", space_id) //@space_id`

	row, err := PGM.pool.Query(context.Background(), sql, pgx.NamedArgs{
		"space_id": space_id,
	})

	defer row.Close()

	if err != nil {
		return 0, "", err
	}

	var Serverid int = -1
	var Owner string

	for row.Next() {
		err = row.Scan(&Serverid, &Owner)
	}

	if Serverid == -1 {
		return Serverid, Owner, errors.New("")
	}

	return Serverid, Owner, nil
}

func (PGM PostgresShareModel) GetServerByShareId(share_id string) (int, error) {

	return 1, nil
}

func (PShareM PostgresShareModel) AddShare(ssr models.SambaShareRequest) (*models.SambaShareResponse, error) {

	server_id, owner_share, err := PShareM.GetServerBySpaceId(ssr.Spaceid)
	if err != nil {
		return nil, err
	}

	if owner_share != ssr.Email {
		return nil, models.ErrorEntryDoesNotExist
	}

	sql := `
		INSERT INTO Samba_Shares (space_id, owner)
		VALUES (@space_id, @owner)
		RETURNING shareid
	`

	rows, err := PShareM.pool.Query(context.Background(), sql, pgx.NamedArgs{
		"space_id": ssr.Spaceid,
		"owner":    ssr.Email,
	})

	if err != nil {
		return nil, err
	}

	var Shareid string
	for rows.Next() {
		err = rows.Scan(&Shareid)
	}

	if err != nil {
		return nil, err
	}
	var index int = -1
	for i, v := range grpc_webclient.GRPCSambaClients {
		if v.Server_id == server_id {
			index = i
			break
		}
	}
	if index == -1 {
		return nil, errors.New("Server Is Not Configured ")
	}

	res, err := grpc_webclient.GRPCSambaClients[index].Grpc_Samba_Client.AllocateSambaShare(context.Background(), &proto_samba_management.RequestSambaShare{
		Owner:    ssr.Email,
		Spaceid:  ssr.Spaceid,
		Password: ssr.Password,
		Shareid:  Shareid,
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
