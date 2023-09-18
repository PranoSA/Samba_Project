package sambaservermanagement_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/PranoSA/samba_share_backend/proto_samba_management"
	sambaservermanagement "github.com/PranoSA/samba_share_backend/samba_server/samba_server_management"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	amqp "github.com/rabbitmq/amqp091-go"
)

func TestCompresS(t *testing.T) {

	//postgresql://postgres:password@127.0.0.1:52269/database_name
	user := "prano"
	password := "prano"
	host := "localhost"
	port := 5432
	database := "samba"
	conn_string := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", user, password, host, port, database)
	fmt.Println(conn_string)

	pool, err := pgxpool.New(context.Background(), conn_string)
	if err != nil {
		t.Error("Failed ToCreate Database connection")
	}
	err = pool.Ping(context.Background())
	if err != nil {
		t.Error("Failed To Ping")
	}

	os.Setenv("SERVER_ID", "1")
	amqp_conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		t.Error("Failed TO Connect TO rabbitmq")
	}

	amqp_conn_server, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		t.Error("Failed TOConnet to RabbitmQ : Server")
	}

	minio_client, err := minio.New("localhost:9000", &minio.Options{
		Creds: credentials.NewStaticV4("minioadmin", "minioadmin", ""),
	})
	if err != nil {
		t.Errorf("Failed To Initialize minio client %v", err)
	}

	var cs sambaservermanagement.CompressServer = sambaservermanagement.CompressServer{
		Rabb_conn:    amqp_conn_server,
		Minio_client: minio_client,
		Pool:         pool,
	}

	ch, err := amqp_conn.Channel()
	if err != nil {
		t.Error("Failed To Initialize AMQP Channel")
	}

	//In This Test, We Are Only Publishing and adding random crap to file share service ....
	err = ch.ExchangeDeclare(proto_samba_management.Exchange_Backup, "direct", true, false, false, false, amqp.Table{})
	if err != nil {
		t.Error("Failed to Declare Exchange")
	}

	queue, err := ch.QueueDeclare(proto_samba_management.Queue_Listening_Backup, true, false, false, false, amqp.Table{})
	if err != nil {
		t.Error("Failed TO Declare QUeue")
	}

	ch.QueueBind(queue.Name, fmt.Sprintf("%s-%d", proto_samba_management.KeyCompressRequest, 1), proto_samba_management.Exchange_Backup, false, amqp.Table{})
	t.Run("publish good point", func(t *testing.T) {

		ch.PublishWithContext(context.Background(), proto_samba_management.Exchange_Backup, proto_samba_management.KeyCompressRequest, false, false, amqp.Publishing{
			Body: []byte("srtrrr"),
		})
		go func() {
			cs.ListenToCompress(1)
		}()
		time.Sleep(5 * time.Second)
	})

}
