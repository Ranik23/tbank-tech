package config

import (
	"log"
	"strings"
	"github.com/spf13/viper"
)

type Config struct {
	ScrapperService 	ScrapperServiceConfig
	DataBase			DataBaseConfig
	Telegram			TelegramConfig 
	TelegramBotServer	TelegramBotServerConfig
	Kafka				KafkaConfig
}


func LoadConfig() (*Config, error) {

	viper.AutomaticEnv()

	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		log.Println("failed to load the .env file")
	}

	return &Config{
		ScrapperService: ScrapperServiceConfig{
			Host: viper.GetString("SCRAPPER_SERVICE_HOST"),
			Port: viper.GetString("SCRAPPER_SERVICE_PORT"),
		},
		DataBase: DataBaseConfig{
			Host: viper.GetString("DATABASE_HOST"),
			Port: viper.GetString("DATABASE_PORT"),
			Username: viper.GetString("DATABASE_USERNAME"),
			Password: viper.GetString("DATABASE_PASSWORD"),
			DBName: viper.GetString("DATABASE_NAME"),
		},
		Telegram: TelegramConfig{
			Token: viper.GetString("TELERGAM_TOKEN"),
		},
		Kafka: KafkaConfig{
			Addresses: strings.Split(viper.GetString("KAFKA_ADDRESSES"), " "),
			Topic: viper.GetString("KAFKA_TOPIC"),
		},
	}, nil


}