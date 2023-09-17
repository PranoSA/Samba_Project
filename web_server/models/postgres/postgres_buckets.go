package postgres_models

import (
	"context"
	"strconv"
	"time"

	"github.com/PranoSA/samba_share_backend/web_server/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

/*
*

	*
	*

CREATE TABLE StreamLinks (

	file_name VARCHAR(256),
	share_id uuid REFERENCES Samba_Shares(shareid),
	PRIMARY KEY(share_id, file_name),
	email VARCHAR(128)

);

CREATE TABLE CompressLinks (

	share_id uuid REFERENCES Samba_Shares(shareid),
	time_backed TIMESTAMP WITH TIME ZONE DEFAULT now(),
	creator VARCHAR(128),
	PRIMARY KEY(share_id, time_backed)

);

/**

);
*/
type PostgresBucketModels struct {
	pool *pgxpool.Pool
}

func (psbm PostgresBucketModels) UserPartOfShare(user string, shareid string) (bool, error) {

	query := `
		SELECT 
		FROM Samba_Share_Users
		where share_id = @share
		AND user_id = @user
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := psbm.pool.Query(ctx, query, &pgx.NamedArgs{
		"share": shareid,
		"user":  user,
	})

	if err == pgx.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil

}

func (psbm PostgresBucketModels) GetCompress(csr models.GetCompressRequests) (*([]models.CompressResponses), error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//Get If User Is Part Of Share

	rows, err := psbm.pool.Query(ctx, `
		SELECT share_id, time_backed
		FROM CompressLinks
		WHERE share_id = @shareid
	`, &pgx.NamedArgs{
		"shareid": csr.Share_id,
	})

	Authorized, _ := psbm.UserPartOfShare(csr.User, csr.Share_id)

	if !Authorized {
		return nil, models.ErrorEntryDoesNotExist
	}

	if err != nil {
		return nil, err
	}

	var QueryRows []models.CompressResponses

	for rows.Next() {
		var newRow models.CompressResponses

		rows.Scan(&newRow.Share_id, &newRow.Timestamp)
		newRow.Url = csr.Share_id + "/" + strconv.FormatInt(newRow.Timestamp.Unix(), 10)
	}

	return &QueryRows, nil
}

func (psbm PostgresBucketModels) CreateCompress(csr models.CreateCompressRequest) (*models.CompressResponses, error) {
	Authorized, _ := psbm.UserPartOfShare(csr.User, csr.Share_id)

	if !Authorized {
		return nil, models.ErrorEntryDoesNotExist
	}

	query := `
		INSERT INTO Compress_Links (user_id, share_id)
		VALUES (@user, @share)
		RETURNING time_backed
	`
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	row := psbm.pool.QueryRow(ctx, query, &pgx.NamedArgs{
		"user":  csr.User,
		"share": csr.Share_id,
	})

	var Time time.Time

	row.Scan(&Time)

	return &models.CompressResponses{
		Share_id:  csr.Share_id,
		Url:       csr.Share_id + "/" + strconv.FormatInt(Time.Unix(), 10),
		Timestamp: Time,
	}, nil

}
