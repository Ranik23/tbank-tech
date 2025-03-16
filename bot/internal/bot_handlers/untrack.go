package bothandlers

import (
	"sync"
	"tbank/bot/internal/service"

	"gopkg.in/telebot.v3"
)


func UnTrackHandler(usecase service.Service, users *sync.Map) telebot.HandlerFunc {
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

		user.state = StateWaitingForLinkUNLINK
		users.Store(userID, user)
		return c.Send("Введите ссылку, которую хотите перестать отслеживать")
	}
}