package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/PranoSA/samba_share_backend/web_server/auth"
	"github.com/PranoSA/samba_share_backend/web_server/controller"
	postgres_models "github.com/PranoSA/samba_share_backend/web_server/models/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
	"gopkg.in/yaml.v3"
)

/**
 *
 * READING IN YAML CONFIG
 *
 * #DEFAULT LOCATION -> /etc/samba_sever/web/web_config.yml

CorsOrigins:
- "localhost:8080"
- "*"

User_Option: "oidc" #"postgres","ldap", "oidc","dynamo"
Session_Option: "oidc" #"oidc", "simple", "redis"
Data_Option: "postgres" #"postgres", "dynamodb", "etcd"

TLS_KEY: ""
TLS_FULLCHAIN: ""


OIDC_CONFIG:
  JWKS_URL: ""
  ISSUER: ""
  AUDIENCE: ""


PG_CONFIG:
  HOST: "localhost"
  PORT: 5432
  USER: "phil"
  USER_TABLE: "users" #THIS IS IGNORED


SAMBA_SERVERS:
- id: 1
  host: localhost ##DEFAULT IS 01.samba_servers.pranoSA
  port: 8080
  CA_CERT: "./ca-cert.pem" #DEFAULT LOCATION?
  TLS_KEY: "./client-key.pem"
  TLS_CERT: "./client-cert.pem"


### Ignored, Only For Example, Remove EXAMPLE_ for actual
EXAMPLE_LDAP_CONFIG:
  TLS_CERT:
  HOST:
  BASE_DN: #BASE DN TO SERACH FROM


EXAMPLE_REDIS_CONFIG:
  CLUSTER_CONNECTION:
    endpoints:
    - host: localhost
      port: 6379
    user: "" #PASSWORD THROUGH CLI
   CA_CERT:
    "./redis-ca-cert.pem"
  CLIENT_CERT:
    "./redis-client-cert.pem"
  CLIENT_KEY:
    "./redis-client-key.pem"


EXAMPLE_DYNAMO_CONFIG:
 *
*/

type ApplicationConfigurations struct {
	https_tls_config *tls.Config
	routes           *controller.AppRouter
	addr             string
	port             int
}

var Application ApplicationConfigurations

type YAMLConfig struct {
	Cors_Origins          []string                    `yaml:"Cors_Origins"`
	OIDC_Config           map[interface{}]interface{} `yaml:"OIDC_CONFIG"`
	User_Config_Option    string                      `yaml:"User_Option"`
	Data_Config_Option    string                      `yaml:"Data_Option"`
	Session_Config_Option string                      `yaml:"Session_Option"`
	TLS_Key               string                      `yaml:"TLS_KEY"`
	Fullchain_Cert        string                      `yaml:"TLS_FULLCHAIN"`
	PG_Config             map[interface{}]interface{} `yaml:"PG_CONFIG"`
	LDAP                  map[interface{}]interface{} `yaml:"LDAP_CONFIG"`
	ETCDConfig            map[interface{}]interface{} `yaml:"ETCD_CONFIG"`
	DynamoDBConfig        map[interface{}]interface{} `yaml:"DYNAMO_CONIG"`
	Redis_Config          map[interface{}]interface{} `yaml:"REDIS_CONFIG"`
}

var ApplicationYamlConfig YAMLConfig

func InitConfig(configPath string) error {
	config_bytes, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(config_bytes, &ApplicationYamlConfig)

	if err != nil {
		return err
	}

	if ApplicationYamlConfig.User_Config_Option != "oidc" && ApplicationYamlConfig.Session_Config_Option == "bearer" {

		return errors.New("For now, Only OIDC can support Bearer Token Authentication")
	}

	if ApplicationYamlConfig.Session_Config_Option == "bearer" && ApplicationYamlConfig.OIDC_Config["jwks_url"].(string) == "" {

	}

	Application.routes.CORS_Origins = ApplicationYamlConfig.Cors_Origins

	if ApplicationYamlConfig.Session_Config_Option == "oidc" {
		auth, err := auth.InitOIDCAuthenticatorFromConfig(ApplicationYamlConfig.OIDC_Config)
		if err != nil {
			return err
		}
		Application.routes.Authenticator = auth

	}

	/**
	 *  Now Here Check For The Other Auth Types ...
	 *  Redis Session
	 *  -> Pass In Backing Store For Users Here ...
	 */

	if ApplicationYamlConfig.Data_Config_Option == "postgres" {
		//Initialize Models Here ...

		conn_string := fmt.Sprintf("%s", ApplicationYamlConfig.PG_Config["Port"].(string))

		pool, err := pgxpool.New(context.Background(), conn_string)
		if err != nil {
			log.Fatal(err)
		}
		Application.routes.Models.Spaces = postgres_models.InitPostgresSpaceModel(pool)
	}

	if ApplicationYamlConfig.Data_Config_Option != "postgres" {
		log.Fatal("Only Postgres Config Implemented \n")
	}

	return nil
}
