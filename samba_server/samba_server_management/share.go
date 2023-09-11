package sambaservermanagement

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/PranoSA/samba_share_backend/proto_samba_management"
	"github.com/minio/minio-go"
	"github.com/minio/minio-go/pkg/credentials"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Need File Path To Read The File

// Need Share ID Because It Needs To Be Part Of Message

//Need Relative FIle Path Because It Needs To Be Part Of message

//It will be in the bucket /shareid/relative_file_path/index.html to reference ....

type MessageTransmitor struct {
	MessageQueue     chan (ShareMessage)
	minioClient      *minio.Client
	rabbit_mq_client *amqp.Connection
}

func InitMessageTransmitor() *MessageTransmitor {
	S3_Endpoint := os.Getenv("S3_ENDPOINT")
	S3_Client_id := os.Getenv("S3_CLIENT_ID")
	S3_Access_Key := os.Getenv("S3_ACCESS_KEY")

	var NewMessageTransmitor MessageTransmitor

	// Initialize minio client object.
	minioClient, err := minio.NewWithOptions(S3_Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(S3_Client_id, S3_Access_Key, ""),
		Secure: false,
	})

	if err != nil {

		log.Fatal(err)
	}

	NewMessageTransmitor.minioClient = minioClient
	NewMessageTransmitor.MessageQueue = make(chan ShareMessage, 100)

	go func() {
		NewMessageTransmitor.ShareMp4()
	}()

	/**
	 *
	 * Declare Exchange and Channel Here to later publish to exchange
	 *
	 */
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	//ch, err := conn.Channel()
	NewMessageTransmitor.rabbit_mq_client = conn

	return &NewMessageTransmitor
}

type ShareMessage struct {
	Filepath           string
	Shareid            string
	relative_file_path string
}

func (mt MessageTransmitor) ShareMp4() error {

	var forever chan struct{}

	ch, err := mt.rabbit_mq_client.Channel()
	if err != nil {

	}

	err = ch.ExchangeDeclare(proto_samba_management.Dash_exchange, "direct", true, false, false, false,
		amqp.Table{})

	if err != nil {
		return err
	}

	queue, err := ch.QueueDeclare(proto_samba_management.Dash_Queue_Requests, true, false, false, false, amqp.Table{})

	err = ch.QueueBind(queue.Name, proto_samba_management.Dash_Request, proto_samba_management.Dash_exchange, false, amqp.Table{})
	if err != nil {
		return err
	}

	for m := range mt.MessageQueue {
		go func(m ShareMessage) {

			_, err := mt.minioClient.FPutObject("backend", m.Shareid+"/"+m.relative_file_path, m.Filepath, minio.PutObjectOptions{})

			if err != nil {
				log.Fatalf("Minio Down For This Service %v", err)
			}

			msg := proto_samba_management.DashMessage{
				Share_id: m.Shareid,
				Filename: m.relative_file_path,
				Resolutions: []struct {
					Width    int
					Height   int
					Bit_Rate int
				}{},
			}

			body_bytes, err := json.Marshal(msg)
			if err != nil {
				forever <- struct{}{}
			}

			ch.PublishWithContext(context.Background(), "samba",
				"dash",
				false, // mandatory
				false, // immediate
				amqp.Publishing{
					ContentType: "application/json",
					Body:        body_bytes,
				})
		}(m)
	}

	<-forever
	return errors.New("Ended")
}
