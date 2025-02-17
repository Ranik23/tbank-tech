package bothandlers

import (
	"sync"
	"tbank/bot/internal/bot-usecase"

	"gopkg.in/telebot.v3"
)


const (
	StateFinished = iota
)

type User struct {
	state   int
	link    string
	tags    []string
	filters []string
}
func HelpHandler(botUseCase botusecase.UseCase, users *sync.Map) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		c.Send("/help /start /list /link /unlink")
		return nil
	}
}