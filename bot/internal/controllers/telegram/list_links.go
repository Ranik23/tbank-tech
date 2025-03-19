package telegram
import (
	"context"
	"fmt"
	"gopkg.in/telebot.v3"
)


func (b *BotHandlers) ListLinksHandler() telebot.HandlerFunc {
	return func(c telebot.Context) error {
		chatID := c.Chat().ID
		responses, err := b.botService.ListLinks(context.Background(), chatID)
		if err != nil {
			return c.Send(fmt.Sprintf("Ошибка: %v", err))
		}

		if len(responses.Links) == 0 {
			c.Send("No Links Tracking")
			return nil
		}

		for _, link := range responses.Links {
			c.Send(link.GetUrl())
		}
		return nil
	}
}