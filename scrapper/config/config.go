package config

import (
	"strings"
	"github.com/spf13/viper"
)




type Config struct {
	ScrapperServer 		ScrapperServerConfig
	DataBase			DataBaseConfig
	BotServer			BotServerConfig
	ScrapperServerHTTP 	ScrapperServerHTTPConfig
	Kafka				KafkaConfig
	BotServerHTTP		BotServerHTTPConfig
	MetricServer		MetricServerConfig
}

func LoadConfig(envPath string) (*Config, error) {

	viper.AutomaticEnv()
	
	viper.SetConfigFile(envPath)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	return &Config{
		ScrapperServer: ScrapperServerConfig{
			Host: viper.GetString("SCRAPPER_SERVICE_HOST"),
			Port: viper.GetString("SCRAPPER_SERVICE_PORT"),
		},
		BotServer: BotServerConfig{
			Host: viper.GetString("TELERGAM_BOT_HOST"),
			Port: viper.GetString("TELEGRAM_BOT_PORT"),
		},
		ScrapperServerHTTP: ScrapperServerHTTPConfig{
			Host: viper.GetString("SCRAPPER_SERVICE_HTTP_HOST"),
			Port: viper.GetString("SCRAPPER_SERVICE_HTTP_PORT"),
		},
		Kafka: KafkaConfig{
			Addresses: strings.Split(viper.GetString("KAFKA_ADDRESSES"), " "),
			Topic: viper.GetString("KAFKA_TOPIC"),
		},
		DataBase: DataBaseConfig{
			Host: viper.GetString("DATABASE_HOST"),
			Port: viper.GetString("DATABASE_PORT"),
			Username: viper.GetString("DATABASE_USERNAME"),
			Password: viper.GetString("DATABASE_PASSWORD"),
			DBName: viper.GetString("DATABASE_NAME"),
			SSL: viper.GetString("DATABASE_SSL"),
		},
		BotServerHTTP: BotServerHTTPConfig{
			Host: viper.GetString("TELEGRAM_BOT_HOST_HTTP"),
			Port: viper.GetString("TELEGRAM_BOT_PORT_HTTP"),
		},
		MetricServer: MetricServerConfig{
			Host: viper.GetString("METRIC_SERVER_HOST"),
			Port: viper.GetString("METRIC_SERVER_PORT"),
		},
	}, nil
}
