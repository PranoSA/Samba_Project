package models

import "errors"

/**
 * Two Connections :
 * One TO Handle SAMBA SERVERS, ISCSI Clients, and
 */

type Models struct {
	Users        UserModel
	Samba_Shares SambaShareModel
	Spaces       SpaceModel
	SambaServers SambaServerModel
}

var (
	ErrorEntryDoesNotExist = errors.New("Entry In DAtabase Does Not Exists")
	ErrorGRPCUnreachable   = errors.New("GRPC Is UnReachable ")
	ErrorNotEnoughSpace    = errors.New("Not Enough SPace")
	ErrorDatabaseTImeout   = errors.New("Database Timeout")
	ErrorMalformedInvite   = errors.New("Malformed Invite")
)
