package sambaservermanagement

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/PranoSA/samba_share_backend/proto_samba_management"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
	amqp "github.com/rabbitmq/amqp091-go"
)

type CompressServer struct {
	Rabb_conn    *amqp.Connection
	Minio_client *minio.Client
	Pool         *pgxpool.Pool
}

func (cs CompressServer) ListenToCompress(id int) {

	/*id := os.Getenv("SERVER_ID")
	if id == "" {
		log.Fatal("Set SERVER_ID Environmental Variable")
	}
	*/

	ch, err := cs.Rabb_conn.Channel()
	if err != nil {
		log.Panic("")
	}

	queue, err := ch.QueueDeclare(proto_samba_management.Queue_Listening_Backup, true, false, false, false, amqp.Table{})

	err = ch.ExchangeDeclare(proto_samba_management.Exchange_Backup, "direct", true, false, false, false, amqp.Table{})

	if err != nil {
		log.Fatalf("Failed To Declare Exchange on Rabbitmq %v", err)
	}

	err = ch.QueueBind(queue.Name, fmt.Sprintf("%s%d", proto_samba_management.KeyCompressRequest, id), proto_samba_management.Exchange_Backup, false, amqp.Table{})

	if err != nil {
		log.Fatalf("Failed To Bind Queue to Exchange on Rabbitmq %v", err)
	}

	forever := make(chan (struct{}), 1)

	msgs, err := ch.Consume(queue.Name, "", true, false, false, false, amqp.Table{})

	go func() {

		for msg := range msgs {

			/**
			 * Decode Json
			 */
			go func(msg amqp.Delivery) {
				var body proto_samba_management.BackupRequest

				err = json.Unmarshal(msg.Body, &body)
				if err != nil {
					fmt.Println(string(msg.Body))
					fmt.Println("Couldn't Unmarshal Queue Message body")
				}

				//Get Directory Where Share Is at
				sharelocation, _ := FindSharePath(cs.Pool, body.Share_id) //UtilFunctionInstance.FindSharePath(body.Share_id)

				var randompathbytes []byte = make([]byte, 24)
				rand.Read(randompathbytes)

				path := base32.StdEncoding.EncodeToString(randompathbytes)

				exec.Command("/usr/bin/tar", "czf", path+".tar.gz", sharelocation)

				timestamp := time.Now().Format(time.RFC3339)

				_, err = cs.Minio_client.FPutObject(context.Background(), proto_samba_management.Bucket_Backup, sharelocation+timestamp, path+".tar.gz", minio.PutObjectOptions{})

				if err != nil {

				}

				/**
				 *
				 * Now Publish To The Complete Queue
				 *
				 */

			}(msg)
		}
	}()

	<-forever
}
