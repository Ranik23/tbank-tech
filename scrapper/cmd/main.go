package main

import (
	"tbank/scrapper/config"
	"tbank/scrapper/internal/app"
)



func main() {

	config, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	app := app.NewApp(config)


	if err := app.Run(); err != nil {
		panic(err)
	}
}