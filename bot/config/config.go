package config

import (
	"log"
	"github.com/spf13/viper"
)


type ScrapperServiceConfig struct {  
	Host string							
	Port string					
}

type TelegramBotServerConfig struct {
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

type TelegramConfig struct {
	Token string
}

type Config struct {
	ScrapperService 	ScrapperServiceConfig
	DataBase			DataBaseConfig
	Telegram			TelegramConfig 
	TelegramBotServer	TelegramBotServerConfig
}


func LoadConfig() (*Config, error) {
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
	}, nil


}