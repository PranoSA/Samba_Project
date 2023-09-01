package models

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
