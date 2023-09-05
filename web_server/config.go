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
			conn_string := fmt.Sprintf("%s", ApplicationYamlConfig.PG_Config["Port"].(string))

			pool, err := pgxpool.New(context.Background(), conn_string)

			if err != nil {
				log.Fatal("")
			}
			sessions.SUO, err = auth.InitPostgresAuth(pool, "brypt")
		}

		Application.routes.Authenticator = sessions
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
