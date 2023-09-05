package models

type SpaceRequest struct {
	Owner     string
	Megabytes int64
}

type SpaceResponse struct {
	Owner     string
	Spaceid   int64
	Email     string
	Megabytes int64
}

type DeleteSpaceRequest struct {
	Owner    string
	Space_id string
}

type SpaceModel interface {
	CreateSpace(SpaceRequest) (*SpaceResponse, error)
	DeleteSpaceById(DeleteSpaceRequest) (*SpaceResponse, error)
	GetSpaceById(DeleteSpaceRequest) (*SpaceResponse, error)
}
