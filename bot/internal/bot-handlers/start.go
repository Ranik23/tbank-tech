package bothandlers

import (
	"context"
	"fmt"
	"sync"
	"tbank/bot/internal/bot-usecase"

	"gopkg.in/telebot.v3"
)


func StartHandler(usecase botusecase.UseCase, users *sync.Map) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		chatID := c.Chat().ID

		if err := usecase.RegisterChat(context.Background(), chatID); err != nil {
			return c.Send(fmt.Sprintf("Ошибка: %v", err))
		}

		c.Send("Чат зарегестрирован!")

		return nil
	}
}