package messagemanager_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

type DashMessage struct {
	Shareid     string
	Resolutions []struct {
		Width  int
		Height int
	}
	Filename string
}

func TestQueue(t *testing.T) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	/*q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")
	*/

	err = ch.ExchangeDeclare(
		"samba",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatal("Failed To Declare Exchange with RabbitMQ Server")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var body []byte

	e := bytes.NewBuffer(body)
	e.Write([]byte{5, 51, 22, 33})
	e1 := make([]byte, 2, 3)
	n, _ := e.Read(e1)
	fmt.Println(n)
	e1 = append(e1, 45)

	body, _ = json.Marshal(DashMessage{
		Shareid: "fart",
		Resolutions: []struct {
			Width  int
			Height int
		}{
			{
				Width:  1000,
				Height: 5000,
			},
			{
				Height: 100,
				Width:  50,
			},
		},
		Filename: "argon",
	})

	/*json.NewEncoder(bytes.NewBuffer(body)).Encode(&DashMessage{
		Shareid: "fart",
		Resolutions: []struct {
			Width  int
			Height int
		}{
			{
				Width:  1000,
				Height: 5000,
			},
			{
				Height: 100,
				Width:  50,
			},
		},
		Filename: "argon",
	})*/

	err = ch.PublishWithContext(ctx,
		"",     // exchange
		"dash", //q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", body)
}

func TestReceiveFunctionality(t *testing.T) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	/*q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")
	*/

	err = ch.ExchangeDeclare(
		"samba",  // name
		"direct", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"dash", // name
		false,  // durable
		false,  // delete when unused
		true,   // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	failOnError(err, "Failed to declare a queue")

	if len(os.Args) < 2 {
		log.Printf("Usage: %s [info] [warning] [error]", os.Args[0])
		os.Exit(0)
	}
	/*for _, s := range os.Args[1:] {
		log.Printf("Binding queue %s to exchange %s with routing key %s",
			q.Name, "logs_direct", s)
		err = ch.QueueBind(
			q.Name,        // queue name
			s,             // routing key
			"logs_direct", // exchange
			false,
			nil)
		failOnError(err, "Failed to bind a queue")
	}
	*/

	err = ch.QueueBind(
		q.Name,  // queue name
		"dash",  // routing key
		"samba", // exchange
		false,
		nil)

	if err != nil {
		t.Error("Failed to Bind")
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)

	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf(" sharted [x] %s", d.Body)
			//d.ContentType
			/*var prettyJSON bytes.Buffer
			error := json.Indent(&prettyJSON, d.Body, "", "\t")
			if error != nil {
				log.Println("JSON parse error: ", error)
				return
			}

			log.Println("CSP Violation:", string(prettyJSON.Bytes()))
			log.Println(string(prettyJSON.Bytes()))
			*/
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
