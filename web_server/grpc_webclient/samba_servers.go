package grpc_webclient

import (
	"log"
	"sync"

	"github.com/PranoSA/samba_share_backend/proto_samba_management"
	"google.golang.org/grpc"
)

type GRPCSambaClient struct {
	Server_id             int
	GRPC_Space_Client     proto_samba_management.SpaceAllocationClient
	Grpc_Samba_Client     proto_samba_management.SambaAllocationClient
	GRPC_Samba_Connection *grpc.ClientConn
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

type GRPCSambaServer struct {
	Host   string
	Ip     string
	Use_IP bool
}

func InitGRPCWebClients(samba []GRPCSambaServer) {

	for i, v := range samba {
		//conn, err := grpc.Dial(samba.Ip)
		conn, err := grpc.Dial(v.Ip)

		if err != nil {
			log.Fatalf("Samba Server %v Can't Be Reached Through GRPC", i)
		}

		client := proto_samba_management.NewSambaAllocationClient(conn)
		client2 := proto_samba_management.NewSpaceAllocationClient(conn)
		GRPCSambaClients = append(GRPCSambaClients, GRPCSambaClient{
			Server_id:             i,
			Grpc_Samba_Client:     client,
			GRPC_Space_Client:     client2,
			GRPC_Samba_Connection: conn,
		})
	}

}

/**
 *
 * How Should We Decide to Split Connections????
 *
 * Belongs To Existing Connection
 *
 *
 */
