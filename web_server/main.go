package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/PranoSA/samba_share_backend/proto_samba_management"
	"github.com/PranoSA/samba_share_backend/web_server/models"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var httpPort string

var https bool
var httpsCertPath string
var httpsKeyPath string

var authMethod string
var sessionMethod string

// ---------------------------------------IF LDAP ---------------------------------//
var ldapURL string //Pass Root Bind Passwords From
var baseDC string  //

// -------------------------------------+ Added Kerberos  ---------------------------//
var LdapTGT string

// -------------------------------- IF POSTGRES ---------------------------------//
var tableName string
var idColumn string
var passwordColumn string

// ---------------------------- IF Bearer -------------------------------//
var jwks_url string

// ------------------------------+ BIND ADDRESS ---------------------------------//
var httpbindAddress string
var httpbindPort int

var SessionType string // OIDC / Bearer, Simple, or Cookie-Based

var configPath string

type Config struct {
	tls_config *tls.Config
	entities   models.Models
	routes     *http.Handler
}

func main() {

	var conn *grpc.ClientConn
	conn, connerr := grpc.Dial(":9000", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if connerr != nil {
		log.Fatalf("did not connect: %s", connerr)
	}
	defer conn.Close()

	c := proto_samba_management.NewSambaAllocationClient(conn)

	response, err := c.AllocateSambaShare(context.Background(), &proto_samba_management.RequestShambaShare{
		Owner:     "pcadler@gmail.com",
		AllocSize: 1000,
	})

	/*response, err := c.SayHello(context.Background(), &chat.Message{Body: "Hello From Client!"}) */
	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}
	log.Printf("Response from server: %s", response.Ip)

	fmt.Printf("Starting Server : \n")

	flag.IntVar(&httpbindPort, "httpport", 80, "Specify The HTTP Port To Listen On By Default = 80")
	flag.StringVar(&httpbindAddress, "address", "127.0.0.1", "Specify Address To Listen On, Default Localhost")
	flag.StringVar(&authMethod, "auth", "oidc", "Specify Auth Methods : Postgres, OIDC, LDAP, RADIUS, KERBEROS+LDAP, PAM")
	flag.StringVar(&sessionMethod, "session", "bearer", "Specify Session Methods : redis, memory, bearer, simpleauth")
	flag.StringVar(&jwks_url, "jwks", "", "Specify a JWK URL if using oidc/bearer session")
	flag.StringVar(&configPath, "config", "/etc/samba_share/web_server/config.yaml", "Config File")

	flag.Parse()

	if configPath == "" {
		configPath = "/etc/samba_server/web/config.yml"
	}

	InitConfig(configPath)

	srv := http.Server{
		Addr:      fmt.Sprintf("%s:%d", Application.addr, Application.port),
		TLSConfig: Application.https_tls_config,
		ErrorLog:  log.Default(),
	}

	log.Fatal(srv.ListenAndServe())
	log.Fatal(srv.ListenAndServeTLS("./", "./"))

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		// handle error!
	}
	defer cli.Close()

}
