package config

import (
	"strings"

	"github.com/spf13/viper"
)

const SERVICENAME = "austin"

type Config struct {
	Server  ServerConfig
	Consul  ConsulConfig
	Storage StorageConfig
}

type ServerConfig struct {
	Host string
	Port int
}

type ConsulConfig struct {
	Hostport    string
	Ttl         string
	Servicename string
}

type StorageConfig struct {
	Host     string
	Port     int
	Keyspace string
}

func GetConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")
	viper.SetConfigType("toml")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	viper.SetEnvPrefix(SERVICENAME)
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()

	cfg := &Config{}
	err = viper.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
