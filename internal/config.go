package internal

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	DBConfig   `yaml:"db"`
	NatsConfig `yaml:"nats"`
}
type DBConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	SSLMode  string `yaml:"sslmode"`
}
type NatsConfig struct {
	Url           string `yaml:"url"`
	StanClusterID string `yaml:"stanClusterID"`
	ClientID      string `yaml:"clientID"`
	Subject       string `yaml:"subject"`
	DurableName   string `yaml:"durableName"`
}

//
//func GetConfig() *Config {
//	return &Config{
//		DBConfig: DBConfig{
//			Host:     "localhost",
//			Port:     "5432",
//			User:     "postgres",
//			Password: "postgres",
//			Database: "l0db",
//			SSLMode:  "disable",
//		},
//		NatsConfig: NatsConfig{
//			Url:     "nats://localhost:4222",
//			Subject: "orders",
//		},
//	}
//}

func MustConfig() *Config {
	viper.SetConfigName("local")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
	config := Config{
		DBConfig: DBConfig{
			Host:     viper.GetString("db.host"),
			Port:     viper.GetString("db.port"),
			User:     viper.GetString("db.user"),
			Password: viper.GetString("db.password"),
			Database: viper.GetString("db.database"),
			SSLMode:  viper.GetString("db.sslmode"),
		},
		NatsConfig: NatsConfig{
			Url:           viper.GetString("nats.url"),
			StanClusterID: viper.GetString("nats.stanClusterID"),
			ClientID:      viper.GetString("nats.clientID"),
			Subject:       viper.GetString("nats.subject"),
			DurableName:   viper.GetString("nats.durableName"),
		},
	}
	return &config
}
