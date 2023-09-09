package main

import (
	"log"
	"net"

	"github.com/PranoSA/samba_share_backend/proto_samba_management"
	sambaservermanagement "github.com/PranoSA/samba_share_backend/samba_server/samba_server_management"
	"google.golang.org/grpc"
)

func main() {

	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	s := sambaservermanagement.SambaServer{}

	proto_samba_management.RegisterSambaAllocationServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}

}
