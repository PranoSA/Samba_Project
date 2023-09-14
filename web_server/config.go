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
	"github.com/PranoSA/samba_share_backend/web_server/grpc_webclient"
	postgres_models "github.com/PranoSA/samba_share_backend/web_server/models/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gopkg.in/yaml.v3"
)

type ApplicationConfigurations struct {
	https_tls_config *tls.Config
	routes           *controller.AppRouter
	addr             string
	port             int
}

var Application ApplicationConfigurations

type YAMLConfig struct {
	Cors_Origins          []string                      `yaml:"Cors_Origins"`
	OIDC_Config           map[interface{}]interface{}   `yaml:"OIDC_CONFIG"`
	User_Config_Option    string                        `yaml:"User_Option"`
	Data_Config_Option    string                        `yaml:"Data_Option"`
	Session_Config_Option string                        `yaml:"Session_Option"`
	TLS_Key               string                        `yaml:"TLS_KEY"`
	Fullchain_Cert        string                        `yaml:"TLS_FULLCHAIN"`
	PG_Config             map[interface{}]interface{}   `yaml:"PG_CONFIG"`
	LDAP                  map[interface{}]interface{}   `yaml:"LDAP_CONFIG"`
	ETCDConfig            map[interface{}]interface{}   `yaml:"ETCD_CONFIG"`
	DynamoDBConfig        map[interface{}]interface{}   `yaml:"DYNAMO_CONIG"`
	Redis_Config          map[interface{}]interface{}   `yaml:"REDIS_CONFIG"`
	Samba_servers         []map[interface{}]interface{} `yaml:"SAMBA_SERVERS"`
	Rabbit_Conf           map[interface{}]interface{}   `yaml:"AMQP_Config"`
}

var ApplicationYamlConfig YAMLConfig

func InitConfig(configPath string) error {
	dir, err := os.Getwd()
	fmt.Println(dir)
	config_bytes, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	Application.routes = &controller.AppRouter{}

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

	if ApplicationYamlConfig.Session_Config_Option == "redis" {
		client := redis.NewClusterClient(&redis.ClusterOptions{
			Addrs: []string{
				"localhost:6379",
			},
			Username: "",
			Password: "",
		})
		sessions := auth.RedisSessionManager{
			RDB: client,
		}

		if ApplicationYamlConfig.User_Config_Option == "postgres" {

			port := ApplicationYamlConfig.PG_Config["PORT"].(int)
			user := ApplicationYamlConfig.PG_Config["USER"].(string)
			host := ApplicationYamlConfig.PG_Config["HOST"].(string)
			database := "samba"
			password := os.Getenv("PG_PASSWORD")

			conn_string := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", user, password, host, port, database)
			//conn_string := fmt.Sprintf("%s", ApplicationYamlConfig.PG_Config["Port"].(string))

			pool, err := pgxpool.New(context.Background(), conn_string)

			if err != nil {
				log.Fatal("")
			}
			sessions.SUO, err = auth.InitPostgresAuth(pool, "brypt")
		}

		Application.routes.Authenticator = sessions
	}

	if ApplicationYamlConfig.Session_Config_Option == "test" {
		os.Setenv("PG_PASSWORD", "prano")
		os.Setenv("RABBITMQ_PASSWORD", "guest")
		Application.routes.Authenticator = auth.AllAllowedAuthenticator{}
	}

	/**
	 *  Now Here Check For The Other Auth Types ...
	 *  Redis Session
	 *  -> Pass In Backing Store For Users Here ...
	 */

	if ApplicationYamlConfig.Data_Config_Option == "postgres" {
		//Initialize Models Here ...

		port := ApplicationYamlConfig.PG_Config["PORT"].(int)
		user := ApplicationYamlConfig.PG_Config["USER"].(string)
		host := ApplicationYamlConfig.PG_Config["HOST"].(string)
		database := "samba"
		password := os.Getenv("PG_PASSWORD")

		conn_string := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", user, password, host, port, database)
		//conn_string := fmt.Sprintf("postgresq://%s:%s@", ApplicationYamlConfig.PG_Config["Port"].(string))

		pool, err := pgxpool.New(context.Background(), conn_string)
		if err != nil {
			log.Fatal(err)
		}
		Application.routes.Models.Spaces = postgres_models.InitPostgresSpaceModel(pool)
		Application.routes.Models.Samba_Shares = postgres_models.InitPostgresShareModel(pool)

		Application.routes.Models.SambaServers = postgres_models.InitPostgresServerModel(pool)
	}

	if ApplicationYamlConfig.Data_Config_Option != "postgres" {
		log.Fatal("Only Postgres Config Implemented \n")
	}

	var servers []grpc_webclient.GRPCSambaServer = make([]grpc_webclient.GRPCSambaServer, len(ApplicationYamlConfig.Samba_servers))
	for i, v := range ApplicationYamlConfig.Samba_servers {
		var next grpc_webclient.GRPCSambaServer = *new(grpc_webclient.GRPCSambaServer)
		next.Id = i
		ip, okip := v["Ip"]

		host, okhost := v["Host"]

		id, okayid := v["ID"]
		if okayid {
			next.Id = id.(int)
		}

		if !okip && !okhost {
			log.Fatalf("Specify Either Ip or Host In config.Samba_Servers[%d].[Ip/Host]", i)
		}

		if okip {
			next.Use_IP = true
			next.Ip = ip.(string)
		}

		if okhost {
			next.Host = host.(string)
		}

		port, portok := v["Port"]
		next.Port = 9887
		if portok {
			next.Port = port.(int)
		}

		servers[i] = next
	}

	rabbitmq_conf := ApplicationYamlConfig.Rabbit_Conf
	user := rabbitmq_conf["RABBIT_USER"].(string)
	password := os.Getenv("RABBITMQ_PASSWORD")
	rabbitport := rabbitmq_conf["PORT"].(int)
	rabbithost := rabbitmq_conf["HOST"].(string)

	rabbitconn, err := amqp091.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d", user, password, rabbithost, rabbitport))
	if err != nil {
		log.Fatal("Failed TO Dial To RabbitMQ Queue")
	}

	Application.routes.Queue = rabbitconn

	grpc_webclient.InitGRPCWebClients(servers)
	return nil
}
