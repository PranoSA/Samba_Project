package controller_test

import (
	"context"
	"errors"

	"github.com/PranoSA/samba_share_backend/proto_samba_management"
	"github.com/PranoSA/samba_share_backend/web_server/models"
)

/*ffffffffff
 *
 *  Sample Interfaces For Evedf
 *
 *
 *
 */

//Implement Everything For Controller Tests Now

type SpaceTestStub struct{} //models.SpaceModel

/**
 * Space Model Data Points
 *
 */
var NextSpaceId string
var FailNext bool

var InitTestSpaceResponses []models.SpaceResponse = []models.SpaceResponse{
	{
		Owner:     "pcadler",
		Spaceid:   "12341234",
		Email:     "pcadler",
		Megabytes: 155,
	},
}

func (stb SpaceTestStub) CreateSpace(sr models.SpaceRequest) (*models.SpaceResponse, error) {
	id := sr.Owner
	megabytes := sr.Megabytes

	if !FailNext {

		InitTestSpaceResponses = append(InitTestSpaceResponses, models.SpaceResponse{
			Spaceid:   NextSpaceId,
			Owner:     id,
			Email:     id,
			Megabytes: megabytes,
		})

		return &models.SpaceResponse{
			Spaceid:   NextSpaceId,
			Owner:     id,
			Email:     id,
			Megabytes: megabytes,
		}, nil
	}
	return nil, errors.New("Failed To Create")
}

func (stb SpaceTestStub) GetSpacesByOwner() {}

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
