package controller_test

import (
	"context"

	"github.com/PranoSA/samba_share_backend/proto_samba_management"
)

/**
 *
 *
 * Use This File To Create a Mock GRPC Server, Pass it In as part of the models
 * on a port
 *
 */

var TestServer GRPCServer

type GRPCServer struct {
	proto_samba_management.UnimplementedSambaAllocationServer
	proto_samba_management.UnimplementedSpaceAllocationServer
}

func (grpcs GRPCServer) AllocateSpace(ctx context.Context, in *proto_samba_management.SpaceAllocationRequest) (*proto_samba_management.SpaceallocationResponse, error) {

	return &proto_samba_management.SpaceallocationResponse{
		Spaceid: in.Spaceid,
	}, nil

}

func (grpcs GRPCServer) AddUserToShare(ctx context.Context, in *proto_samba_management.AddUser) (*proto_samba_management.AddUserResponse, error) {
	return &proto_samba_management.AddUserResponse{}, nil
}
