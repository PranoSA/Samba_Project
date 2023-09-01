package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/PranoSA/samba_share_backend/web_server/auth"
	"github.com/PranoSA/samba_share_backend/web_server/models"
	clientv3 "go.etcd.io/etcd/client/v3"
	"gopkg.in/yaml.v3"
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

	fmt.Printf("Starting Server : \n")

	flag.IntVar(&httpbindPort, "httpport", 80, "Specify The HTTP Port To Listen On By Default = 80")
	flag.StringVar(&httpbindAddress, "address", "127.0.0.1", "Specify Address To Listen On, Default Localhost")
	flag.StringVar(&authMethod, "auth", "oidc", "Specify Auth Methods : Postgres, OIDC, LDAP, RADIUS, KERBEROS+LDAP, PAM")
	flag.StringVar(&sessionMethod, "session", "bearer", "Specify Session Methods : redis, memory, bearer, simpleauth")
	flag.StringVar(&jwks_url, "jwks", "", "Specify a JWK URL if using oidc/bearer session")
	flag.StringVar(&configPath, "config", "/etc/samba_share/web_server/config.yaml", "Config File")

	flag.Parse()

	auth.InitOIDCAuthenticator(jwks_url, []string{}, []string{})
	fmt.Printf(configPath + "\n")

	config_bytes, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)

		fmt.Errorf("Failed To Read File %v : %v", configPath, err)
	}

	type YAMLConfig struct {
		Aud                    []string                    `yaml:"acceptedAud"`
		Iss                    []string                    `yaml:"acceptedIss"`
		PG_Config              map[interface{}]interface{} `yaml:"DB_URI"`
		Jwks                   string                      `yaml:"JWKS_URI"`
		LDAP                   map[interface{}]interface{} `yaml:"LDAP_CONFIG"`
		ServerManagementMethod string                      `yaml:"Server_Management"`
		ETCDConfig             map[interface{}]interface{} `yaml:"ETCD_Config"`
	}

	server_conf := YAMLConfig{}

	err = yaml.Unmarshal(config_bytes, &server_conf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(server_conf.Aud[0])

	var problems int = 0

	if authMethod != "oidc" && sessionMethod == "bearer" {
		problems++
		fmt.Errorf("Bearer Only Works If OIDC Is Specified")
	}

	if sessionMethod == "bearer" && jwks_url == "" {
		problems++
		fmt.Errorf("If using Session Method of Bearer Must include a json wek key url")
	}

	/***
	 * ETCD CONFIGURATION IF EXISTS
	 *
	 */

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		// handle error!
	}
	defer cli.Close()

}
