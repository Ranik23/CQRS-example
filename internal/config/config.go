package config

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Storage     StorageConfig  `mapstructure:"storage"`
	Server      ServerConfig   `mapstructure:"server"`
	Logging     LoggingConfig  `mapstructure:"logging"`
	Kafka 	    KafkaConfig    `mapstructure:"kafka"`
	Redis       RedisConfig    `mapstructure:"redis"`
}

func LoadConfig(envPath string, configPath string) (*Config, error) {
	err := godotenv.Load(envPath)
	if err != nil {
		fmt.Println("Error loading .env file")
		return nil, err
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
		return nil, err
	}

	viper.AutomaticEnv() // Read environment variables

	viper.BindEnv("storage.main.host", "MAIN_HOST")
	viper.BindEnv("storage.main.password", "MAIN_PASSWORD")
	viper.BindEnv("storage.main.username", "MAIN_USERNAME")
	viper.BindEnv("storage.main.dbname", "MAIN_DBNAME")
	viper.BindEnv("storage.main.port", "MAIN_PORT")

	viper.BindEnv("storage.side.host", "SIDE_HOST")
	viper.BindEnv("storage.side.password", "SIDE_PASSWORD")
	viper.BindEnv("storage.side.username", "SIDE_USERNAME")
	viper.BindEnv("storage.side.dbname", "SIDE_DBNAME")
	viper.BindEnv("storage.side.port", "SIDE_PORT")

	viper.BindEnv("server.port", "SERVER_PORT")
	viper.BindEnv("kafka.brokers", "KAFKA_BROKERS")
	viper.BindEnv("redis.address", "REDIS_ADDRESS")

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %v", err)
	}

	log.Printf("Loaded configuration: %+v", config)

	return &config, nil
}
