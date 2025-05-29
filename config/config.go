package config

import (
	"log"
	"os"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

type Config struct {
	DB       DBConfig
	JWT      JWTConfig
	Server   ServerCfg
	InfluxDB InfluxDBConfig
	Proxy    ProxyConfig
	Ditto    DittoConfig
}

type DittoConfig struct {
	URL      string `envconfig:"DITTO_URL"`
	Username string `envconfig:"DITTO_USERNAME"`
	Password string `envconfig:"DITTO_PASSWORD"`
	WSURL    string `envconfig:"DITTO_WS_URL"`
}

type DBConfig struct {
	Host               string `envconfig:"DB_HOST" default:"localhost"`
	Port               string `envconfig:"DB_PORT" default:"5432"`
	User               string `envconfig:"DB_USER" default:"postgres"`
	Password           string `envconfig:"DB_PASSWORD" default:"postgres"`
	DBName             string `envconfig:"DB_NAME" default:"postgres"`
	SSLMode            string `envconfig:"SSL_MODE" default:"disable"`
	SetMaxIdleConns    string `envconfig:"SET_MAX_IDLE_CONNS" default:""`
	SetMaxOpenConns    string `envconfig:"SET_MAX_OPEN_CONNS" default:""`
	SetConnMaxLifetime string `envconfig:"SET_CONN_MAX_LIFETIME" default:""`
}

type JWTConfig struct {
	Secret                string `envconfig:"JWT_SECRET"`
	ExpirationTime        string `envconfig:"JWT_EXPIRATION_TIME"`
	RefreshSecret         string `envconfig:"JWT_REFRESH_SECRET"`
	RefreshExpirationTime string `envconfig:"JWT_REFRESH_EXPIRATION_TIME"`
}

type ServerCfg struct {
	ServerURL  string `envconfig:"SERVER_URL" default:"localhost"`
	Port       string `envconfig:"PORT" default:"3001"`
	Env        string `envconfig:"ENVIRONMENT" default:"development"`
	GINMode    string `envconfig:"GIN_MODE" default:"debug"`
	Production bool   `envconfig:"PRODUCTION" default:"false"`
}

type InfluxDBConfig struct {
	URL    string `envconfig:"INFLUXDB_URL"`
	Token  string `envconfig:"INFLUXDB_TOKEN"`
	Org    string `envconfig:"INFLUXDB_ORG"`
	Bucket string `envconfig:"INFLUXDB_BUCKET"`
}

type ProxyConfig struct {
	AuthUsername string `envconfig:"PROXY_AUTH_USERNAME"`
	AuthPassword string `envconfig:"PROXY_AUTH_PASSWORD"`
	TargetURL    string `envconfig:"PROXY_TARGET_URL"`
	WSURL        string `envconfig:"PROXY_WS_URL"`
}

func NewConfig() (*Config, error) {
	LoadConfig()

	var cfg Config

	if err := envconfig.Process("", &cfg.Ditto); err != nil {
		log.Fatalf("Failed to process Ditto config: %v", err)
	}

	if err := envconfig.Process("", &cfg.DB); err != nil {
		log.Fatalf("Failed to process DB config: %v", err)
	}
	if err := envconfig.Process("", &cfg.JWT); err != nil {
		log.Fatalf("Failed to process JWT config: %v", err)
	}
	if err := envconfig.Process("", &cfg.Server); err != nil {
		log.Fatalf("Failed to process Server config: %v", err)
	}
	if err := envconfig.Process("", &cfg.InfluxDB); err != nil {
		log.Fatalf("Failed to process InfluxDB config: %v", err)
	}
	if err := envconfig.Process("", &cfg.Proxy); err != nil {
		log.Fatalf("Failed to process Proxy config: %v", err)
	}

	return &cfg, nil
}

func LoadConfig() {
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	for _, env := range viper.AllKeys() {
		if viper.GetString(env) != "" {
			_ = os.Setenv(env, viper.GetString(env))
			_ = os.Setenv(strings.ToUpper(env), viper.GetString(env))
		}
	}
}

var Module = fx.Options(
	fx.Provide(NewConfig),
)
