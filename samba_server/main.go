package main

import (
	"log"
	"net"

	"github.com/PranoSA/samba_share_backend/proto_samba_management"
	"google.golang.org/grpc"
)

func main() {

	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	s := SambaServer{}

	proto_samba_management.RegisterSambaAllocationServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
