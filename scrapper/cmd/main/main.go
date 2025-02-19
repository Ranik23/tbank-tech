package main

import (
	"tbank/scrapper/internal/app"
)



func main() {

	app, err := app.NewApp()
	if err != nil {
		panic(err)
	}

	if err := app.Run(); err != nil {
		panic(err)
	}
}