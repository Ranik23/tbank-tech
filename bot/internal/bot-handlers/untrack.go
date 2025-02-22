package bothandlers

import (
	"sync"
	"tbank/bot/internal/bot-usecase"

	"gopkg.in/telebot.v3"
)


func UnTrackHandler(usecase botusecase.UseCase, users *sync.Map) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		userID := c.Sender().ID

		userRaw, exists := users.Load(userID)
		var user *User
		if !exists {
			return c.Send("Вам нужно зарегестрироваться. Используйте команду /start")
		} 
		
		user = userRaw.(*User)
		user.state = StateWaitingForLinkUNLINK
		users.Store(userID, user)
		return c.Send("Введите ссылку, которую хотите перестать отслеживать")
	}
}