package bothandlers

import (
	"context"
	"fmt"
	"sync"
	botusecase "tbank/bot/internal/bot_usecase"

	"gopkg.in/telebot.v3"
)


func ListHandler(usecase botusecase.UseCase, users *sync.Map) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		chatID := c.Chat().ID
		
		responses, err := usecase.ListLinks(context.Background(), chatID)
		if err != nil {
			return c.Send(fmt.Sprintf("Ошибка: %v", err))
		}
		for _, link := range responses.Links {
			c.Send(link)
		}
		return nil
	}
}