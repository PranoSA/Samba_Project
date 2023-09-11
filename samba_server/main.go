package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/PranoSA/samba_share_backend/proto_samba_management"
	sambaservermanagement "github.com/PranoSA/samba_share_backend/samba_server/samba_server_management"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

/**
 *
 * Transition to Adding DYnamoDB LAter ...
 */

var Port int
var ServerId int

func main() {

	flag.IntVar(&Port, "port", 9887, "Port For GRPC Server To Listen To")
	flag.IntVar(&ServerId, "serverid", -1, "Id Of Server")

	flag.Parse()

	if ServerId == -1 {
		log.Fatal("Please Specify Server ID")
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	connstring := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", "prano", "prano", "localhost", 5432, "samba")

	pool, err := pgxpool.New(context.Background(), connstring)
	if err != nil {
		log.Fatalf("Failed To Connect To Pool %v", err)
	}

	err = pool.Ping(context.TODO())
	if err != nil {
		log.Fatalf("Failed To Ping %v", err)
	}

	grpcServer := grpc.NewServer()

	//s := sambaservermanagement.SambaServer{}
	s := sambaservermanagement.NewSambaServer(pool, ServerId)

	proto_samba_management.RegisterSambaAllocationServer(grpcServer, s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}

}
