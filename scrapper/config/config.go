package config

import (
	"strings"

	"github.com/spf13/viper"
)


type KafkaConfig struct {
	Addresses []string
}

type ScrapperServerConfig struct {
	Host string
	Port string
}


type ScrapperServerHTTPConfig struct {
	Host string
	Port string
}

type DataBaseConfig struct {		
	Host 		string				
	Port 		string			
	Username 	string				
	Password 	string					
	DBName 		string					
}

type BotServerConfig struct {
	Host string
	Port string
}


type Config struct {
	ScrapperServer 		ScrapperServerConfig
	DataBase			DataBaseConfig
	Bot					BotServerConfig
	ScrapperServerHTTP 	ScrapperServerHTTPConfig
	Kafka				KafkaConfig
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	return &Config{
		ScrapperServer: ScrapperServerConfig{
			Host: viper.GetString("SCRAPPER_SERVICE_HOST"),
			Port: viper.GetString("SCRAPPER_SERVICE_PORT"),
		},
		Bot: BotServerConfig{
			Host: viper.GetString("TELERGAM_BOT_HOST"),
			Port: viper.GetString("TELEGRAM_BOT_PORT"),
		},
		ScrapperServerHTTP: ScrapperServerHTTPConfig{
			Host: viper.GetString("SCRAPPER_SERVICE_HTTP_HOST"),
			Port: viper.GetString("SCRAPPER_SERVICE_HTTP_PORT"),
		},
		Kafka: KafkaConfig{
			Addresses: strings.Split(viper.GetString("KAFKA_ADDRESSES"), " "),
		},
	}, nil
}