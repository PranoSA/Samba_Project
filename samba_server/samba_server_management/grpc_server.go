package sambaservermanagement

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

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

	/**
	 *
	 * Here We Want to DO The SQL Query Instead of On the PG Server
	 *
	 */

	fs := FS.ChooseOne(in.Size)

	sql := `
		INSERT INTO Samba_Spaces (fs_id, owner, alloc_size, time_created)
		VALUES (@fs_id, @owner, @size, @time_created)
		RETURNING fs_id, owner
	`

	row := s.pool.QueryRow(context.Background(), sql, pgx.NamedArgs{
		"fs_id":        fs.Fsid,
		"owner":        in.Owner,
		"size":         in.Size,
		"time_created": time.Now(),
	})

	/*if row. {
		return &proto_samba_management.SpaceallocationResponse{
			StatusCode: 1,
		}, nil
	}

	if fs == nil {
		return &proto_samba_management.SpaceallocationResponse{
			Spaceid:    "",
			StatusCode: 2,
		}, nil
	}
	*/

	var fs_id string
	var e error
	row.Scan(&fs_id, &e)

	if fs_id == "" {
		return &proto_samba_management.SpaceallocationResponse{
			StatusCode: 1,
		}, nil
	}
	_ = os.Mkdir(fs.MouthPath+"/"+fs_id, 0771)

	return &proto_samba_management.SpaceallocationResponse{
		Spaceid:    in.Spaceid,
		StatusCode: 0,
		Size:       in.Size,
		Fsid:       fs_id,
	}, nil
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
	SELECT fsid, device, mnt_point, capacity -  COALESCE((
		SELECT SUM (Samba_Spaces.alloc_size)
		FROM Samba_Spaces
		WHERE Samba_Spaces.fs_id = Samba_File_Systems.fsid
		),0) AS space_left
		FROM Samba_File_Systems
		WHERE server_id = @server;
	`

	/**
	 * 	SELECT fsid, device, COALESCE(mnt_point,''), Samba_File_Systems.capacity - SUM(Samba_Spaces.alloc_size)
	FROM Samba_File_Systems
	JOIN Samba_Spaces
	ON Samba_Spaces.fs_id = Samba_File_Systems.fsid
	WHERE Samba_File_Systems.server_id = 1
	GROUP BY Samba_File_Systems.fsid
	 *
	*/

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
	SELECT spaceid, device, COALESCE(mnt_point, ''), fs_id 
	FROM Samba_Spaces
	JOIN Samba_File_Systems
	ON Samba_File_Systems.fsid = Samba_Spaces.fs_id
	WHERE spaceid = @space_id
	`

	row, err := s.pool.Query(context.Background(), sql, pgx.NamedArgs{
		"space_id": space_id,
	})

	if err != nil {
		return "", " "
	}

	var Spaceid, Dev, Mount_path, Fs_id string
	for row.Next() {
		err = row.Scan(&Spaceid, &Dev, &Mount_path, &Fs_id)
		if err != nil {

		}
	}

	if Mount_path == "" {
		return fmt.Sprintf("/mount/samba_server/%s/%s", Fs_id, space_id), Fs_id
	}

	return Mount_path, Fs_id

}

func (s *SambaServer) AllocateSambaShare(ctx context.Context, in *proto_samba_management.RequestSambaShare) (*proto_samba_management.SambaResponse, error) {

	/**
	 *
	 * Find Where Space is and Create New Folder for Spaceid
	 *
	 */
	fmt.Println("Got Request")
	path, fsid := s.FindSpacePath(in.Spaceid) //Find Space Path, Now What ...????
	fmt.Println(path)
	CreateSambaShare(path, in.Shareid, in.Owner, in.Password, in.Spaceid)
	//CreateSambaShare(path, in.Owner, in.Spaceid, in.Shareid, in.Password)

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
	 * To Call The Script I Need the ...
	 * Shareid and Spaceid:::
	 *
	 */

	user := in.User
	password := in.Password
	shareid := in.ShareId

	row, err := s.pool.Query(context.Background(), "SELECT space_id FROM Samba_Shares where shareid=@shareid", pgx.NamedArgs{
		"shareid": shareid,
	})

	if err != nil {

	}

	var spaceid string
	row.Scan(&spaceid)

	AddUserToShareId(user, password, shareid, spaceid)

	return &proto_samba_management.AddUserResponse{}, nil
}
