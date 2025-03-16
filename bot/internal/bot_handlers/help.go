package bothandlers

import (
	"sync"
	"tbank/bot/internal/service"

	"gopkg.in/telebot.v3"
)

type User struct {
	state   int
	link    string
	tags    []string
	filters []string
}

func HelpHandler(usecase service.Service, users *sync.Map) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		c.Send("/help /start /list /track /untrack")
		return nil
	}
}