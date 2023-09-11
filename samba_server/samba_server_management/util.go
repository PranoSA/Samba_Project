package sambaservermanagement

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UtilFunction struct {
	pool *pgxpool.Pool
}

var UtilFunctionInstance UtilFunction

func (uf UtilFunction) FindSpacePath(space_id string) (string, string) {

	sql := `
	SELECT spaceid, dev, COALESCE(mount_path, ''), fs_id 
	FROM Samba_Spaces
	WHERE spaceid = @space_id
	`

	row, err := uf.pool.Query(context.Background(), sql, pgx.NamedArgs{
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

func (uf UtilFunction) FindSharePath(shareid string) (string, error) {

	sql := `
		SELECT space_id
		FROM Samba_Shares
		WHERE share_id = @shareid
	`

	row, err := uf.pool.Query(context.Background(), sql, pgx.NamedArgs{
		"share_id": shareid,
	})
	if err != nil {
		return "", err
	}

	var spaceid string
	row.Scan(&spaceid)

	space_path, _ := uf.FindSpacePath(spaceid)

	return space_path + "/" + shareid, nil
}

func FindSharePath(pool *pgxpool.Pool, shareid string) (string, error) {

	sql := `
		SELECT space_id
		FROM Samba_Shares
		WHERE share_id = @shareid
	`

	row, err := pool.Query(context.Background(), sql, pgx.NamedArgs{
		"share_id": shareid,
	})
	if err != nil {
		return "", err
	}

	var spaceid string
	row.Scan(&spaceid)

	space_path, _ := FindSpacePath(pool, spaceid)

	return space_path + "/" + shareid, nil
}

func FindSpacePath(pool *pgxpool.Pool, space_id string) (string, string) {

	sql := `
	SELECT spaceid, dev, COALESCE(mount_path, ''), fs_id 
	FROM Samba_Spaces
	WHERE spaceid = @space_id
	`

	row, err := pool.Query(context.Background(), sql, pgx.NamedArgs{
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
