package messagemanager

import (
	"context"
	"log"

	"github.com/minio/minio-go/v7"
	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageManager struct {
	S3_Client *minio.Client
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

type Job struct {
	Share_id  string //Samba Share FIles Belongs To
	File_path string //File Path To Serve From
}

func (ms MessageManager) PublicDash(share_id string, file_path string) {
	/**
	 *
	 * Listen To Queue And Publish To Endpoint
	 *
	 * Read From -> bucket : /severs -> share_id/filepath
	 *
	 * Publish To -> bucket : public ->  /share_id/filepath/[.m3u8] -> and all the associated manifests
	 *
	 * FilePath -> Random ID That The Script Will Write To
	 * -> CALL FFMEPG ON THAT FILE
	 * -> Rewrite All THe Files to a BUCKET IN MINIO..?????
	 */

	err := ms.S3_Client.FGetObject(context.Background(), "backend", share_id+"/"+file_path, "", minio.GetObjectOptions{})

	if err != nil {

	}
	//Do Work ....

	ms.S3_Client.FPutObject(context.Background(), "frontend", share_id+"/"+file_path, "", minio.PutObjectOptions{})

}

func StartListen() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
