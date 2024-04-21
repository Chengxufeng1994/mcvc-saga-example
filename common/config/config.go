package config

import (
	"github.com/ThreeDotsLabs/watermill"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ApplicationConfig struct {
	LogConfig      Log            `mapstructure:"log"`
	PostgresConfig PostgresConfig `mapstructure:"postgres"`
	RedisConfig    RedisConfig    `mapstructure:"redis"`
	NatsConfig     NatsConfig     `mapstructure:"nats"`
	JaegerConfig   JaegerConfig   `mapstructure:"jaeger"`
	JWTConfig      JWTConfig      `mapstructure:"jwt"`
	RpcEndpoints   RpcEndpoints   `mapstructure:"rpc_endpoints"`
}

type Log struct {
	Level string `mapstructure:"level"`
}

type PostgresConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DbName       string `mapstructure:"db_name"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
}

type RedisDeploymentType string

const (
	Cluster    RedisDeploymentType = "cluster"
	Standalone RedisDeploymentType = "standalone"
)

type RedisConfig struct {
	Type        RedisDeploymentType `mapstructure:"type"`
	Addrs       string              `mapstructure:"addrs"`
	User        string              `mapstructure:"user"`
	Password    string              `mapstructure:"password"`
	Db          int                 `mapstructure:"db"`
	PoolSize    int                 `mapstructure:"pool_size"`
	PoolTimeout int                 `mapstructure:"pool_timeout"`
	ReadOnly    bool                `mapstructure:"read_only"`
	MaxRetries  int                 `mapstructure:"max_retries"`
	Subscriber  Subscriber          `mapstructure:"subscriber"`
}

type Subscriber struct {
	ConsumerID    string `mapstructure:"consumer_id"`
	ConsumerGroup string `mapstructure:"consumer_group"`
}

type JaegerConfig struct {
	Endpoint string `mapstructure:"endpoint"`
}

type NatsConfig struct {
	Host           string          `mapstructure:"host"`
	Port           int             `mapstructure:"port"`
	ClusterID      string          `mapstructure:"cluster_id"`
	ClientID       string          `mapstructure:"client_id"`
	NatsSubscriber *NatsSubscriber `mapstructure:"subscriber"`
}

type NatsSubscriber struct {
	QueueGroup  string `mapstructure:"queue_group"`
	DurableName string `mapstructure:"durable_name"`
}

type JWTConfig struct {
	Secret              string `mapstructure:"secret"`
	AccessTokenExpires  int    `mapstructure:"access_token_expires"`
	RefreshTokenExpires int    `mapstructure:"refresh_token_expires"`
}

type RpcEndpoints struct {
	AuthServiceHost    string `mapstructure:"auth_service_host"`
	ProductServiceHost string `mapstructure:"product_service_host"`
}

func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, DisableColors: false, TimestampFormat: "2006-01-02 15:04:05"})
}

func LoadApplicationConfig(path string) *ApplicationConfig {
	log.Infoln("common.Config initializing...")
	var config *ApplicationConfig
	viper.AddConfigPath(path)
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		panic(err)
	}

	if config.RedisConfig.Subscriber.ConsumerID == "" {
		config.RedisConfig.Subscriber.ConsumerID = watermill.NewShortUUID()
	}

	return config
}
