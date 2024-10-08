package configs

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConfig `mapstructure:"server"`
	DB     DBConfig     `mapstructure:"db"`
	NATS   NATSConfig   `mapstructure:"nats"`
	Kafka  KafkaConfig  `mapstructure:"kafka"`
	SubscriberConfig
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

type DBConfig struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"ssl_mode"`
}

type NATSConfig struct {
	Host        string `mapstructure:"host"`
	Port        string `mapstructure:"port"`
	ClusterID   string `mapstructure:"cluster_id"`
	ClientID    string `mapstructure:"client_id"`
	DurableName string `mapstructure:"durable_name"`
	SubjectPost string `mapstructure:"subject_post"`
}

type SubscriberConfig struct {
	DurableName string
	SubjectPost string
}

type KafkaConfig struct {
	Brokers []string `mapstructure:"brokers"`
	Topic   string   `mapstructure:"topic"`
	GroupID string   `mapstructure:"group_id"`
}

func LoadConfig(path string) (Config, error) {
	var config Config
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		return config, err
	}
	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}
	return config, nil
}
