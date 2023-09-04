package grpc_webclient

import (
	"sync"

	"github.com/PranoSA/samba_share_backend/proto_samba_management"
)

type GRPCSambaClient struct {
	Server_id         int
	GRPC_Space_Client proto_samba_management.SpaceAllocationClient
	Grpc_Samba_Client proto_samba_management.SambaAllocationClient
}

var Next_id int = 0
var IncrMutex sync.Mutex = sync.Mutex{}
var GRPCSambaClients []GRPCSambaClient

/**
 * TODO: Add More Complicated Load Balancing Later
 */
func GetAndUpdateNextId() int {

	IncrMutex.Lock()
	Prev_id := Next_id
	Next_id = (Next_id + 1) % len(GRPCSambaClients)
	IncrMutex.Unlock()
	return Prev_id
}

/**
 *
 * How Should We Decide to Split Connections????
 *
 * Belongs To Existing Connection
 *
 *
 */
