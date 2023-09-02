package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/PranoSA/samba_share_backend/proto_samba_management"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FileSystem struct {
	Fsid      string
	Dev       string
	RoomLeft  int64
	MouthPath string
}

type FileSystems struct {
	FileSystems []FileSystem
}

var FS FileSystems

func (f *FileSystems) ChooseOne(capacity int64) *FileSystem {

	for _, fs := range f.FileSystems {
		if fs.RoomLeft > capacity {

			fs.RoomLeft = fs.RoomLeft - capacity
			return &fs
		}
	}
	return nil
}

type SambaServer struct {
	*proto_samba_management.UnimplementedSambaAllocationServer
	pool *pgxpool.Pool
}

type Disk struct {
	fsid       string
	server_id  int
	device     string
	mnt_point  *string
	capacity   int
	space_left int
}

func NewSambaServer(pool *pgxpool.Pool, serverid int) *SambaServer {

	sql := `
	SELECT fsid, device, mnt_point, Samba_File_Systems.capacity - SUM(Samba_Spaces.alloc_size)
	FROM Samba_File_Systems
	JOIN Samba_Spaces
	ON Samba_Spaces.fs_id = Samba_File_Systems.fsid
	WHERE Samba_File_Systems.server_id = @server
	GROUP BY Samba_File_Systems.fsid
	`

	disksrows, query_err := pool.Query(context.Background(), sql, pgx.NamedArgs{
		"server": serverid,
	})

	defer disksrows.Close()

	if query_err != nil {
		log.Fatal("Failed To Initialize Connection To Postgres")
	}

	var NextDisk FileSystem = FileSystem{}
	for disksrows.Next() {
		if disksrows.Err() != nil {
			log.Fatal("Failed To Initialize Connection To Postgres")
		}

		err := disksrows.Scan(&NextDisk.Fsid, &NextDisk.Dev, &NextDisk.MouthPath, &NextDisk.RoomLeft)
		if err != nil {
			log.Fatal("Failed To Initialize Connection To Postgres")
		}

		FS.FileSystems = append(FS.FileSystems, NextDisk)
	}

	return &SambaServer{pool: pool}
}

func (s *SambaServer) AllocateSambaShare(ctx context.Context, in *proto_samba_management.RequestShambaShare) (*proto_samba_management.SambaResponse, error) {
	fmt.Println("Got Request")

	f := FS.ChooseOne(in.AllocSize)

	if f == nil {
		return &proto_samba_management.SambaResponse{
			Status: 1,
			Fsid:   "",
			Ip:     "",
		}, errors.New("Couldn't ALlocate FS")
	}

	return &proto_samba_management.SambaResponse{}, nil
}
