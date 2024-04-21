package bootstrap

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type BootstrapConfig struct {
	Application string `mapstructure:"application"`
	Environment string `mapstructure:"environment"`
	GinMode     string `mapstructure:"gin_mode"`
	HTTP        HTTP   `mapstructrue:"http"`
	Grpc        Grpc   `mapstructure:"grpc"`
}

type HTTP struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	ReadTimeout  string `mapstructure:"read_timeout"`
	WriteTimeout string `mapstructure:"write_timeout"`
	IdleTimeout  string `mapstructure:"idle_timeout"`
}

type Grpc struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// Config bootstrap configuration
var Config *BootstrapConfig

func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, DisableColors: false, TimestampFormat: "2006-01-02 15:04:05"})
}

func LoadBootstrapConfig(path string) *BootstrapConfig {
	log.Infoln("common.Bootstrap initializing...")
	var config *BootstrapConfig
	viper.AddConfigPath(path)
	viper.AddConfigPath(".")
	viper.SetConfigName("bootstrap")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		panic(err)
	}

	return config
}
