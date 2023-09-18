package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresServer struct {
	pool *pgxpool.Pool
}

func (ps PostgresServer) AddServer(id int, ip string) (int, int, error) {
	sql := `
	INSERT INTO samba_server (serverid, lastip)
	VALUES (@id, @ip)
	RETURNING *
`

	row, err := ps.pool.Query(
		context.Background(),
		sql,
		pgx.NamedArgs{
			"id": id,
			"ip": ip,
		},
	)

	defer row.Close()

	if err != nil {
		log.Fatalf("ID %v Could Possibly Exist, %v", id, err)
		return 0, 0, err
	}

	var return_id int
	var return_ipd int

	row.Scan(&return_id, &return_ipd)

	if row.Next() {
		log.Fatal("Error In Program : Adding Duplicated Rows")
	}

	return return_id, return_ipd, nil

}

/**
* SCHEMA:
	CREATE TABLE Samba_File_Systems (
*     fsid uuid PRIMARY KEY DEFAULT gen_random_uuid(),
   server_id INTEGER NOT NULL REFERENCES Samba_Server(serverid),
   device VARCHAR(128),
   mnt_point VARCHAR(255),
   capacity INTEGER NOT NULL
*
*/

type Disk struct {
	fsid       string
	server_id  int
	device     string
	mnt_point  *string
	capacity   int
	space_left int
}

/**
 * SCHEMA:
 *
 * CREATE TABLE Samba_Spaces (
    spaceid uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    fs_id uuid NOT NULL REFERENCES Samba_File_Systems(fsid),
    owner VARCHAR(128),
    alloc_size INTEGER NOT NULL
);
 *
*/

func (ps PostgresServer) GetDisksOnServer(serverid int) ([]Disk, error) {

	/*sql := `
		SELECT fsid, server_id, device, mnt_point, capacity, Samba_File_Systems.capacity - SUM(Samba_Spaces.alloc_size)
		FROM Samba_File_Systems
		JOIN Samba_Spaces
		ON Samba_Spaces.fs_id = Samba_File_Systems.fsid
		WHERE Samba_File_Systems.server_id = @server
		GROUP BY Samba_File_Systems.fsid
	`*/

	sql := `
	SELECT fsid, server_id, device, mnt_point, capacity, capacity -  COALESCE((
		SELECT SUM (Samba_Spaces.alloc_size)
		FROM Samba_Spaces
		WHERE Samba_Spaces.fs_id = Samba_File_Systems.fsid
		),0) AS space_left
		FROM Samba_File_Systems
		WHERE server_id = @server;
		`

	disksrows, query_err := ps.pool.Query(context.Background(), sql, pgx.NamedArgs{
		"server": serverid,
	})

	defer disksrows.Close()

	if query_err != nil {
		return []Disk{}, query_err
	}

	var QueryDisks []Disk

	var NextDisk Disk = Disk{}

	for disksrows.Next() {
		if disksrows.Err() != nil {
			return []Disk{}, disksrows.Err()
		}

		err := disksrows.Scan(&NextDisk.fsid, &NextDisk.server_id, &NextDisk.device, &NextDisk.mnt_point, &NextDisk.capacity, &NextDisk.space_left)
		if err != nil {
			return []Disk{}, err
		}

		QueryDisks = append(QueryDisks, NextDisk)
	}

	return QueryDisks, nil
}

func (ps PostgresServer) AddDiskToServer(serverid int, device string, mount_path string, capacity int) error {

	sql := `
		INSERT INTO Samba_File_Systems (server_id, device, mnt_point, capacity)
		VALUES (@id, @device, @mount_path, @capacity)
		RETURNING fsid
	`

	row := ps.pool.QueryRow(context.Background(), sql, pgx.NamedArgs{
		"id":       serverid,
		"disk":     device,
		"device":   device,
		"capacity": capacity,
	})

	/*defer row.Close()

	if queryerr != nil {
		return queryerr
	}
	*/

	var fsid string

	scanerr := row.Scan(&fsid)
	if scanerr != nil {
		return scanerr
	}

	/*mount_dir := fmt.Sprintf("/samba_app/mnt/%d/%s", serverid, fsid)

	mkdirerr := os.Mkdir(mount_dir, 0770)
	if mkdirerr != nil {
		return mkdirerr
	}
	*/

	/**
	 * Need To Add Cross - REFERENCE HERE
	 */

	fmt.Printf("Added FS : %s \n", fsid)

	return nil

}

/**
 *
 * 	FOR ADMINISTRATION PURPOSE ONLY
 *
 * !!! DO NOT DO THIS MANUALLY FOR END USERS, THIS IS ONLY FOR TESTING FUNCTIONALITY
 */
