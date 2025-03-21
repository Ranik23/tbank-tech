package telegram

import (
	"context"
	"github.com/Ranik23/tbank-tech/bot/internal/models"
	"gopkg.in/telebot.v3"
)

func (b *BotHandlers) MessageHandler() telebot.HandlerFunc {
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

		text := c.Text()

		switch user.State {
		case StateWaitingForLinkUNLINK:
			user.Link = text
			user.State = StateFinished
			b.users.Store(userID, user)

			response, err := b.botService.RemoveLink(context.Background(), userID, user.Link)
			if err != nil {
				return c.Send("❌ Ошибка удаления ссылки!")
			}
			return c.Send(response.GetMessage(), telebot.ModeHTML)
		case StateWaitingForLinkLINK:
			user.Link = text
			user.State = StateFinished
			b.users.Store(userID, user)

			response, err := b.botService.AddLink(context.Background(), userID, user.Link)
			if err != nil {
				return c.Send("❌ Ошибка добавления ссылки!")
			}
			return c.Send(response.GetMessage(), telebot.ModeHTML)
		case StateWaitingForTheToken:
			user.Token = text
			if user.Name == "" {
				return c.Send("⚠️ Ошибка: имя пользователя не установлено. Попробуйте заново.", telebot.ModeHTML)
			}

			response, err := b.botService.RegisterUser(context.Background(), userID, user.Name, user.Token)
			if err != nil {
				user.State = StateWaitingForTheName
				return c.Send("❌ Ошибка регистрации пользователя!")
			}
			return c.Send(response.GetMessage(), telebot.ModeHTML)
		case StateWaitingForTheName:
			user.Name = text
			user.State = StateWaitingForTheToken
			b.users.Store(userID, user)
			return c.Send("🔑 Введите персональный токен:")
		}
		return nil
	}
}
