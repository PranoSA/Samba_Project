package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"time"

	psa "github.com/PranoSA/samba_share_backend/proto_samba_management"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

/**
 *
 *	Enforce Insertion With Pre-Allocated Fs Identifiers
 *
 */

type FileSystemRequest struct {
	Fs_id      string
	Mount_path string
	Server_id  int
	Size       int64
	Dev        string
}

func NewFileSystemRequest(pool *pgxpool.Pool, fs FileSystemRequest, config *tls.Config) error {

	var conn *grpc.ClientConn
	//Replce with Actual DNS Later
	conn, connerr := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if connerr != nil {
		log.Fatalf("did not connect: %s", connerr)
		return connerr
	}
	defer conn.Close()

	tx, err := pool.BeginTx(context.Background(), pgx.TxOptions{})

	if err != nil {
		return err
	}

	insert_statement := `
	INSERT INTO Samba_File_Systems(server_id, device, mnt_point, capacity)
	VALUES(@serverid, @dev, @mount, @capacity)
	RETURNING fsid
	`
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	row := tx.QueryRow(ctx, insert_statement, &pgx.NamedArgs{
		"serverid": fs.Server_id,
		"dev":      fs.Dev,
		"mount":    fs.Mount_path,
		"capacity": fs.Size,
	})
	if err != nil {
		tx.Rollback(context.Background())
		return err
	}

	var Fsid string
	err = row.Scan(&Fsid)
	if err != nil {
		tx.Rollback(context.Background())
		fmt.Println(err)
		return err
	}

	c := psa.NewDiskAllocationClient(conn)

	res, err := c.AddDiskToServer(context.Background(), &psa.PartitionAllocRequest{
		Device:    fs.Dev,
		MountPath: fs.Mount_path,
		AllocSize: fs.Size,
		Fsid:      Fsid,
	})

	fmt.Println(res.StatusMessage)

	if err != nil {
		tx.Rollback(context.Background())
		return err
	}

	if res.StatusCode != 0 {
		tx.Rollback(context.Background())
		return errors.New("Failed To Allocate Disk, Recheck Proper Mount Points")
	}

	tx.Commit(context.Background())

	return nil
}
