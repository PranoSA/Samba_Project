package models

type SambaShareRequest struct {
	Email    string
	Spaceid  string
	Password string
}

type SambaShareResponse struct {
	Email   string
	Shareid string
}

type ShareInviteRequest struct {
	Email   string
	Shareid string
}

type ShareInviteResponse struct {
	Email       string
	Inviteid    string
	Invite_code string
}

type ShareInviteAccept struct {
	Email       string
	Inviteid    string
	Invite_code string
	Password    string
}

type SambaShareModel interface {
	AddShare(SambaShareResponse) (*SambaShareResponse, error)
	DeleteShare(SambaShareResponse) (*SambaShareResponse, error)
	CreateInvite(ShareInviteRequest) (*ShareInviteResponse, error)
	AcceptInvite(ShareInviteAccept) (*ShareInviteResponse, error)
}
