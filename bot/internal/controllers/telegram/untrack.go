package telegram

import (
	"github.com/Ranik23/tbank-tech/bot/internal/models"

	"gopkg.in/telebot.v3"
)

func (b *BotHandlers) UnTrackHandler() telebot.HandlerFunc {
	return func(c telebot.Context) error {
		userID := c.Sender().ID

		userRaw, exists := b.users.Load(userID)
		var user *models.User
		if !exists {
			user = &models.User{State: StateFinished}
			b.users.Store(userID, user)
		} else {
			user = userRaw.(*models.User)
		}

		user.State = StateWaitingForLinkUNLINK
		b.users.Store(userID, user)
		return c.Send("Введите ссылку, которую хотите перестать отслеживать")
	}
}
