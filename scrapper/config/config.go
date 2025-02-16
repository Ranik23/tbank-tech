package config

import "github.com/spf13/viper"


type ScrapperServerConfig struct {
	Host string
	Port string
}


type Config struct {
	ScrapperServer ScrapperServerConfig
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile("env")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	return &Config{
		ScrapperServer: ScrapperServerConfig{
			Host: viper.GetString("SCRAPPER_SERVICE_HOST"),
			Port: viper.GetString("SCRAPPER_SERVICE_PORT"),
		},
	}, nil
}