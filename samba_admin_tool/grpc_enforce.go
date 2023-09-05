package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"

	psa "github.com/PranoSA/samba_share_backend/proto_samba_management"
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
	conn, connerr := grpc.Dial("localhost:9000", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if connerr != nil {
		log.Fatalf("did not connect: %s", connerr)
	}
	defer conn.Close()

	c := psa.NewDiskAllocationClient(conn)

	res, err := c.AddDiskToServer(context.Background(), &psa.PartitionAllocRequest{
		Device:    fs.Dev,
		MountPath: fs.Mount_path,
		AllocSize: fs.Size,
		Fsid:      fs.Fs_id,
	})

	if err != nil {
	}

	fmt.Println(res)

	return nil
}
