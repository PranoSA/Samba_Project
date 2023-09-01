package models

import (
	"github.com/jackc/pgx/v5/pgxpool"
	clientv3 "go.etcd.io/etcd/client/v3"
)

/**
 *
 * This file is for load balancing between samba servers and finding backend severs
 *
 * For Example
 *
 * You get a request for a share, you call a backend samba server and allocate that share
 * and track samba shares against backend servers
 * You have to figure out where this backend server is
 *
 *
 *
 *
 */

type SambaServerModel interface {
}

type SambaServerETCD struct {
	db *clientv3.Client
}

type SambaServerPostgres struct {
	client *pgxpool.Pool
}
