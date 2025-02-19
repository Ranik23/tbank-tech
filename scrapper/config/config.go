package config

import "github.com/spf13/viper"


type ScrapperServerConfig struct {
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
	}, nil
}