package sambaservermanagement

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/PranoSA/samba_share_backend/proto_samba_management"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FileSystem struct {
	Fsid           string
	Dev            string
	RoomLeft       int64
	MouthPath      string
	Lock           *sync.Mutex
	Transaction_id int32
}

type FileSystems struct {
	FileSystems []FileSystem
}

var FS FileSystems

func (f *FileSystems) ChooseOne(capacity int64) *FileSystem {

	for _, fs := range f.FileSystems {
		fs.Lock.Lock()
		if fs.RoomLeft > capacity {
			fs.RoomLeft = fs.RoomLeft - capacity
			return &fs
		}
		fs.Lock.Unlock()

		InitFromDiskLabels(f.FileSystems)
	}
	return nil
}

type SambaServer struct {
	*proto_samba_management.UnimplementedSambaAllocationServer
	*proto_samba_management.UnimplementedDiskAllocationServer
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

func (s SambaServer) AlloateSpace(ctx context.Context, in *proto_samba_management.SpaceAllocationRequest) (*proto_samba_management.SpaceallocationResponse, error) {

	return &proto_samba_management.SpaceallocationResponse{}, nil
}

func (s SambaServer) AllocateSpaceConversation(stream proto_samba_management.SambaAllocation_AllocateSpaceConversationServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		/**
		 *
		 * What We Want To Do Here is Define A start_time, and anything
		 * We Want to Return
		 */
		if in.Sequence == 1 {
			server := FS.ChooseOne(in.Size)
			server.Lock.Lock()
			defer server.Lock.Unlock()
			// server.Transaction_id = server.Transaction_id

		}

		if in.Sequence == 2 {

		}

	}
}

func NewSambaServer(pool *pgxpool.Pool, serverid int) *SambaServer {

	sql := `
	SELECT fsid, device, COALESCE(mnt_point,''), Samba_File_Systems.capacity - SUM(Samba_Spaces.alloc_size)
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
	NextDisk.Lock = &sync.Mutex{}

	var ErrorSystems []FileSystem = []FileSystem{}

	for disksrows.Next() {
		if disksrows.Err() != nil {
			log.Fatal("Failed To Initialize Connection To Postgres")
		}

		err := disksrows.Scan(&NextDisk.Fsid, &NextDisk.Dev, &NextDisk.MouthPath, &NextDisk.RoomLeft)

		if err != nil {
			log.Fatal("Failed To Initialize Connection To Postgres")
		}

		var test_path string

		test_path = NextDisk.MouthPath
		if NextDisk.MouthPath == "" {
			test_path = fmt.Sprintf("/mnt/samba_server/%d/%s", serverid, NextDisk.Fsid)
		}

		mode, err := os.Stat(test_path)

		if err != nil {
			ErrorSystems = append(ErrorSystems, NextDisk)
			fmt.Printf("Expected File System Mounted at %s\n", test_path)
		}

		if !mode.IsDir() {
			ErrorSystems = append(ErrorSystems, NextDisk)
			fmt.Printf("%s is a file not a directory\n", test_path)
		}

		if mode.Mode()&0700 != 0700 {
			fmt.Printf("Improper Permissions On %s\n", test_path)
			ErrorSystems = append(ErrorSystems, NextDisk)
		}

		FS.FileSystems = append(FS.FileSystems, NextDisk)

	}

	for _, fs := range ErrorSystems {
		test_path := fs.MouthPath
		if fs.MouthPath == "" {
			test_path = fmt.Sprintf("/mnt/samba_server/%d/%s", serverid, fs.Fsid)
		}
		fmt.Printf("Expected %s moutned at %s \n", fs.Dev, test_path)
	}

	if len(ErrorSystems) > 0 {
		log.Fatal("")
	}

	return &SambaServer{pool: pool}
}

/**
 * 	This is For if fs_id is already pre-generated and mounted
 */

func (s *SambaServer) AddDiskToServer(ctx context.Context, in *proto_samba_management.PartitionAllocRequest) (*proto_samba_management.PartitionAllocResponse, error) {

	//Check Multiple Things

	//Space On Disk Is More than Allocated

	//Accessible On Mount Point or

	mount_path := in.MountPath
	if mount_path == "" {
		mount_path = fmt.Sprintf("/mount/samba_server/%v", in.Fsid)
	}

	ok, err := EnsureMount(in.Device, mount_path)
	if !ok {
		fmt.Printf("%v", err)
		return &proto_samba_management.PartitionAllocResponse{
			StatusCode:    1,
			StatusMessage: "Wrong Mount Point For Disk",
		}, err
	}

	return &proto_samba_management.PartitionAllocResponse{
		StatusCode: 0,

		StatusMessage: "",
	}, nil
}

func (s *SambaServer) FindSpacePath(space_id string) (string, string) {

	sql := `
	SELECT spaceid, dev, COALESCE(mount_path, ''), fs_id 
	FROM Samba_Spaces
	WHERE spaceid = @space_id
	`

	row, err := s.pool.Query(context.Background(), sql, pgx.NamedArgs{
		"space_id": space_id,
	})

	if err != nil {
		return "", " "
	}

	var spaceid, dev, mount_path, fs_id string
	row.Scan(&spaceid, &dev, &mount_path, &fs_id)

	if mount_path == "" {
		return fmt.Sprintf("/mount/samba_server/%s/%s", fs_id, space_id), fs_id
	}

	return mount_path, fs_id

}

func (s *SambaServer) AllocateSambaShare(ctx context.Context, in *proto_samba_management.RequestShambaShare) (*proto_samba_management.SambaResponse, error) {

	/**
	 *
	 * Find Where Space is and Create New Folder for Spaceid
	 *
	 */

	_, fsid := s.FindSpacePath(in.Spaceid) //Find Space Path, Now What ...????

	return &proto_samba_management.SambaResponse{
		Status: 0,
		Fsid:   fsid,
		Ip:     "",
	}, nil
}

func (s *SambaServer) AddUserToShare(ctx context.Context, in *proto_samba_management.AddUser) (*proto_samba_management.AddUserResponse, error) {

	/**
	 *
	 * Call Commands To Create the Share And Samba Formatting
	 *
	 *
	 *
	 */

	return &proto_samba_management.AddUserResponse{}, nil
}
