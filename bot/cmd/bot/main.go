package main

import (
	"github.com/Ranik23/tbank-tech/bot/internal/app"
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
