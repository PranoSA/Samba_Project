package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/PranoSA/samba_share_backend/proto_samba_management"
	"github.com/julienschmidt/httprouter"
	"github.com/rabbitmq/amqp091-go"
)

/**
 *
 * Eventually Add These Features
 *
 */

type DashHTTPRequest struct {
	Share_id    string
	File_name   string
	Resolutions []struct {
		Width   int
		Height  int
		Bitrate int
	}
}

type CompressRequest struct {
	Share_id  string
	File_name string
}

var DashPublications chan (DashHTTPRequest) = make(chan DashHTTPRequest, 10000)
var CompressPublications chan (DashHTTPRequest) = make(chan DashHTTPRequest, 10000)

func (ap AppRouter) StartDashPublisher() {
	go func() {
		forever := make(chan struct{})

		ch, err := ap.Queue.Channel()
		if err != nil {
			log.Fatal(err)
		}

		ch.ExchangeDeclare(proto_samba_management.Dash_exchange, "direct", true, false, false, false, amqp091.Table{})

		for m := range DashPublications {
			serverid, _ := ap.Models.SambaServers.GetServerByShareId(m.Share_id)
			ch.PublishWithContext(context.Background(), proto_samba_management.Dash_exchange, fmt.Sprintf("%s%d", proto_samba_management.Dash_Request, serverid), false, false, amqp091.Publishing{})
		}

		<-forever
	}()
}

func (ap AppRouter) StartCompressPublisher() {
	go func() {
		forever := make(chan struct{})

		ch, err := ap.Queue.Channel()
		if err != nil {
			log.Fatal(err)
		}

		ch.ExchangeDeclare(proto_samba_management.Exchange_Backup, "direct", true, false, false, false, amqp091.Table{})

		for m := range DashPublications {
			serverid, _ := ap.Models.SambaServers.GetServerByShareId(m.Share_id)
			ch.PublishWithContext(context.Background(), proto_samba_management.Exchange_Backup, fmt.Sprintf("%s%d", proto_samba_management.Queue_Listening_Backup, serverid), false, false, amqp091.Publishing{})
		}

		<-forever
	}()
}

func (ar AppRouter) RequestDash(w http.ResponseWriter, r *http.Request, ap httprouter.Params) {

	email := r.Context().Value("Authorization")

	var request DashHTTPRequest

	json.NewDecoder(r.Body).Decode(&request)

	DashPublications <- request

	fmt.Println(email)
}

func (ar AppRouter) CompressShare(w http.ResponseWriter, r *http.Request, ap httprouter.Params) {

}

func (ar AppRouter) GetCompressLinks(w http.ResponseWriter, r *http.Request, ap httprouter.Params) {

}
