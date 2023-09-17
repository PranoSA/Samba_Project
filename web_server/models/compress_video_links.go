package models

import "time"

type DashHTTPRequest struct {
	Share_id    string
	File_name   string
	Resolutions []struct {
		Width   int
		Height  int
		Bitrate int
	}
}

type DashHTTPResponse struct {
	Share_id  string
	File_name string
	URL       string
}

type GetVideoRequests struct {
	Share_id string
	User     string
}

type getVideoResponse struct {
}

type GetCompressRequests struct {
	Share_id string
	User     string
}

type CompressResponses struct {
	Share_id  string
	Url       string
	Timestamp time.Time
}

type CreateCompressRequest struct {
	Share_id string
	User     string
}

type VideoModels interface {
	CreateVideoRequest(DashHTTPRequest) (*DashHTTPResponse, error)
	GetVideoRequests(GetVideoRequests) (*([]DashHTTPResponse), error)
}

type CompressModels interface {
	GetCompress(GetCompressRequests) (*([]CompressResponses), error)
	CreateCompress(CreateCompressRequest) (*CompressResponses, error)
}
