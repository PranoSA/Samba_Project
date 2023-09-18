package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/PranoSA/samba_share_backend/proto_samba_management"
	sambaservermanagement "github.com/PranoSA/samba_share_backend/samba_server/samba_server_management"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
)

/**
 *
 * Transition to Adding DYnamoDB LAter ...
 */

var Port int
var ServerId int

var MinioPort int
var MinioHost string
var MinioUser string
var MinioPassword string

var RabbitMQHost string
var RabbitMQPort int
var RabbitMQUser string
var RabbitMQPassword string

var CompressionServer sambaservermanagement.CompressServer

func main() {

	flag.IntVar(&Port, "port", 9887, "Port For GRPC Server To Listen To")
	flag.IntVar(&ServerId, "serverid", -1, "Id Of Server")

	flag.IntVar(&RabbitMQPort, "amqpport", 5672, "Port To Listen To RabbitMQ For")
	flag.StringVar(&RabbitMQHost, "amqphost", "localhost", "Host of RabbitMQ")
	flag.StringVar(&RabbitMQUser, "amqpuser", "guest", "")
	flag.StringVar(&RabbitMQPassword, "amqppassword", "guest", "")

	flag.IntVar(&MinioPort, "minioport", 9000, "Port To Listen To RabbitMQ For")
	flag.StringVar(&MinioHost, "miniohost", "localhost", "Host of RabbitMQ")
	flag.StringVar(&MinioUser, "miniouser", "minioadmin", "")
	flag.StringVar(&MinioPassword, "miniopassword", "minioadmin", "")

	flag.Parse()

	if ServerId == -1 {
		log.Fatal("Please Specify Server ID")
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	connstring := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", "prano", "prano", "localhost", 5432, "samba")

	pool, err := pgxpool.New(context.Background(), connstring)
	if err != nil {
		log.Fatalf("Failed To Connect To Pool %v", err)
	}

	err = pool.Ping(context.TODO())
	if err != nil {
		log.Fatalf("Failed To Ping %v", err)
	}

	grpcServer := grpc.NewServer()

	//s := sambaservermanagement.SambaServer{}
	s := sambaservermanagement.NewSambaServer(pool, ServerId)

	amqp_conn_string := fmt.Sprintf("amqp://%s:%s@%s:%d", RabbitMQUser, RabbitMQPassword, RabbitMQHost, RabbitMQPort)
	amqp_conn, err := amqp.Dial(amqp_conn_string)
	if err != nil {
		log.Fatal("Failed TO Connect TO rabbitmq")
	}

	/*minio_client, err := minio.New(fmt.Sprintf("%s:%d", MinioHost, MinioPort), &minio.Options{
		Creds: credentials.NewStaticV4(MinioUser, MinioPassword, ""),
	})
	*/
	minio_conn := fmt.Sprintf("%s:%d", MinioHost, MinioPort)

	minio_client, err := minio.New(minio_conn, &minio.Options{
		Creds: credentials.NewStaticV4(MinioUser, MinioPassword, ""),
	})

	if err != nil {
		log.Fatalf("Failed To Initialize minio client %v", err)
	}

	if err != nil {
		log.Fatalf("Failed To Initialize minio client %v", err)
	}

	CompressionServer.Pool = pool

	CompressionServer.Rabb_conn = amqp_conn

	CompressionServer.Minio_client = minio_client

	//Set Up Minio Client Here
	//os.Setenv("SERVER_ID", ServerId)
	CompressionServer.ListenToCompress(ServerId)

	proto_samba_management.RegisterSambaAllocationServer(grpcServer, s)
	proto_samba_management.RegisterDiskAllocationServer(grpcServer, s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}

}
