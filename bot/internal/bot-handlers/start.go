package bothandlers

import (
	"sync"
	"tbank/bot/internal/bot-usecase"

	"gopkg.in/telebot.v3"
)


func StartHandler(usecase botusecase.UseCase, users *sync.Map) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		userID := c.Sender().ID

		userRaw, exists := users.Load(userID)
		var user *User
		if !exists {
			user = &User{state: StateFinished}
			users.Store(userID, user)
		} else {
			user = userRaw.(*User)
		}

		user.state = StateWaitingForTheToken
		users.Store(userID, user)

		return c.Send("Введите персональный токен")
	}
}