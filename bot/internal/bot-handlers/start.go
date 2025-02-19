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

		response, err := usecase.RegisterChat(context.Background(), chatID);
		if err != nil {
			return c.Send(fmt.Sprintf("Ошибка: %v", err))
		}
		c.Send(response.GetMessage())

		return nil
	}
}