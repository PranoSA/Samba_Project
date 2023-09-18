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

//var _ models.SpaceModel = SpaceTestStub{}

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

func (stb SpaceTestStub) GetSpacesByOwner(sr models.SpaceResponse) ([]models.SpaceResponse, error) {
	var spaces []models.SpaceResponse

	for _, s := range InitTestSpaceResponses {
		if s.Email == sr.Owner {
			spaces = append(spaces, s)
		}
	}

	return spaces, nil
}

func (stb SpaceTestStub) DeleteSpaceById(dsr models.DeleteSpaceRequest) (*models.SpaceResponse, error) {

	var Response models.SpaceResponse
	var found int = 0

	for i, _ := range InitTestSpaceResponses {
		if InitTestSpaceResponses[i].Spaceid == dsr.Space_id {
			if InitTestSpaceResponses[i].Email != dsr.Owner {
				return nil, models.ErrorEntryDoesNotExist
			}
			Response = InitTestSpaceResponses[i]
			InitTestSpaceResponses = append(InitTestSpaceResponses[0:i], InitTestSpaceResponses[i+1:]...)
			found = i
		}
	}
	if found == 0 {
		return nil, models.ErrorEntryDoesNotExist
	}
	return &Response, nil
}

//func (stm SpaceTestStub) GetSpaceById(in models.DeleteSpaceRequest) (*models.SpaceResponse, error) {}

//func (stb SpaceTestStub)

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
