package models

type SpaceRequest struct {
	email     string
	megabytes int64
}

type SpaceResponse struct {
	spaceid   int64
	email     string
	megabytes int64
}

type SpaceModel interface {
	CreateSpace(SpaceRequest) (SpaceResponse, error)
	DeleteSpaceById(int64) (SpaceResponse, error)
	GetSpaceById(int64) (SpaceResponse, error)
}
