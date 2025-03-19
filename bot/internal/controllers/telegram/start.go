package telegram

import (
	"github.com/Ranik23/tbank-tech/bot/internal/models"

	"gopkg.in/telebot.v3"
)

func (b *BotHandlers) StartHandler() telebot.HandlerFunc {
	return func(c telebot.Context) error {
		userID := c.Sender().ID

		userRaw, exists := b.users.Load(userID)
		var user *models.User
		if !exists {
			user = &models.User{}
			b.users.Store(userID, user)
		} else {
			user = userRaw.(*models.User)
		}

		user.State = StateWaitingForTheName
		b.users.Store(userID, user)

		return c.Send("📝 Введите свое имя:", telebot.ModeHTML)
	}
}
