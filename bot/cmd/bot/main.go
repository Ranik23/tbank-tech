package main

import (
	"tbank/bot/internal/bot/handlers"
	"tbank/bot/internal/bot/router"
	"tbank/bot/config"
	"time"
	"gopkg.in/telebot.v3"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	bot, err := telebot.NewBot(telebot.Settings{
		Token: cfg.Telegram.Token,
		Poller: &telebot.LongPoller{
			Timeout: 10 * time.Second,
			AllowedUpdates: []string{
				"message",
				"edited_message",
			},
		},
	})
	if err != nil {
		panic(err)
	}

	router := router.NewRouter(bot, nil)

	router.AddHandler("/start", handlers.NewStartHandler(nil))
	router.AddHandler("/help", handlers.NewHelpHandler(nil))
	router.AddHandler("/track", handlers.NewTrackHandler(nil))
	router.AddHandler("/untrack", handlers.NewUntrackHandler(nil))
	router.AddHandler("/list", handlers.NewListHandler(nil))

	router.RegisterHandlers()

	bot.Start()

}
