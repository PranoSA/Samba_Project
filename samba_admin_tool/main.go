package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Port int
var Host string
var Password string
var User string
var Verbose bool

type ProcessConfig struct {
	manager PostgresServer
}

var Config ProcessConfig

var Command bool

func main() {
	args := os.Args

	flag.BoolVar(&Command, "AddServer", false, "")
	flag.BoolVar(&Command, "AddDisk", false, "")
	flag.BoolVar(&Command, "DiskUtil", false, "")
	flag.BoolVar(&Command, "SambaAdd", false, "")

	cmd := args[1]

	if cmd == "--AddServer" {
		var id int
		var ip string

		flag.StringVar(&Host, "host", "localhost", "Host Of The Postgres Database Server")
		flag.IntVar(&Port, "port", 5432, "Port of the Postgres Database Server")
		flag.StringVar(&User, "user", "postgres", "Postgres User")
		flag.StringVar(&Password, "password", "postgres", "Password for User")
		flag.BoolVar(&Verbose, "v", false, "Print Connection Information")

		flag.IntVar(&id, "id", 0, "ID Of Server To Add")
		flag.StringVar(&ip, "ip", "", "Ip of Server to Add")

		flag.Parse()

		if id == 0 {
			fmt.Printf("Please Enter Valid Server ID\n")
			os.Exit(0)
		}

		if ip == "" {
			fmt.Printf("Please Enter Valid Server IP Address\n")
			os.Exit(0)
		}

		fmt.Println(User)

		pool, err := pgxpool.New(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%d/samba?sslmode=verify-ca&pool_max_conns=10", User, Password, Host, Port))

		if Verbose {
			fmt.Printf("Connection:%v \n ", pool.Config().ConnString())
		}

		if err != nil {
			log.Fatalf("%v", err)
		}

		err = pool.Ping(context.Background())

		if err != nil {
			log.Fatalf("%v", err)
		}

		Config.manager.pool = pool

		Config.manager.AddServer(id, ip)

	}

	if cmd == "--DiskUtil" {
		var id int
		var ip string
		var device string
		var remote bool
		var mount_path string
		var size int

		var view bool
		var Enforce bool

		flag.StringVar(&Host, "host", "localhost", "Host Of The Postgres Database Server")
		flag.IntVar(&Port, "port", 5432, "Port of the Postgres Database Server")
		flag.StringVar(&User, "user", "postgres", "Postgres User")
		flag.StringVar(&Password, "password", "postgres", "Password for User")
		flag.BoolVar(&Verbose, "v", false, "Print Connection Information")

		flag.StringVar(&device, "disk", "", "Block Device To Add To Samba Server")
		flag.BoolVar(&remote, "remote", false, "Remote Device? ")

		flag.BoolVar(&view, "view", false, "Only Print Contents Of Mounted Disks")

		flag.IntVar(&id, "id", 0, "ID Of Server To Add")
		flag.StringVar(&ip, "ip", "", "Ip of Server to Add")
		flag.StringVar(&mount_path, "mount", "", "Mount Location")
		flag.BoolVar(&Enforce, "enforce", false, "Enforce GRPC Communication with The Server")

		flag.IntVar(&size, "size", 0, "Size of Disk To Be Added")

		flag.Parse()

		if id == 0 {
			fmt.Printf("Please Enter Valid Server ID\n")
			os.Exit(0)
		}

		if ip == "" {
			fmt.Printf("Please Enter Valid Server IP Address\n")
			os.Exit(0)
		}

		pool, err := pgxpool.New(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%d/samba?sslmode=verify-ca&pool_max_conns=10", User, Password, Host, Port))

		if Verbose {
			fmt.Printf("Connection:%v \n ", pool.Config().ConnString())
		}

		if err != nil {
			log.Fatalf("%v", err)
		}

		err = pool.Ping(context.Background())

		if err != nil {
			log.Fatalf("%v", err)
		}

		Config.manager.pool = pool

		disks, err := Config.manager.GetDisksOnServer(id)

		if err != nil {
			log.Fatalf("Failed %v\n", err)
		}

		var diskAlreadyExists bool = false

		for _, disk := range disks {
			fmt.Printf("Disk %v: \n", disk.fsid)
			fmt.Printf("\t Server Id : %v \n", disk.server_id)
			fmt.Printf("\t Device : %v \n", disk.device)
			if disk.device == device {
				diskAlreadyExists = true
			}
			fmt.Printf("\t Mount_Path : %v \n", disk.mnt_point)
			fmt.Printf("\t Capacity : %v \n", disk.capacity)
			fmt.Printf("\t Space Left : %v \n", disk.space_left)
			fmt.Println("---------------------------------------------------------------- ")
			fmt.Println("")
		}

		fmt.Printf("%v Disks On Server \n", len(disks))

		if view {
			return
		}

		if diskAlreadyExists {
			fmt.Println("Disk Already Exists On Server")
			os.Exit(0)
		}

		if device == "" {
			fmt.Println("Please Enter Valid Disk")
			os.Exit(0)
		}

		if mount_path == "" {
			mount_location := fmt.Sprintf("/samba_app/mnt/%d/fsid", id)

			fmt.Printf("Mounting at Default Location of %v \n", mount_location)
		}

		if size == 0 {
			fmt.Println("Please Enter a Valid Disk Size with --size")
			os.Exit(0)
		}

		err = Config.manager.AddDiskToServer(id, device, mount_path, size)

		if err != nil {
			log.Fatalf("%v \n", err)
		}
	}

}
