package telegram
import (
	"context"
	"gopkg.in/telebot.v3"
)


func (b *BotHandlers) ListLinksHandler() telebot.HandlerFunc {
	return func(c telebot.Context) error {
		userID := c.Sender().ID
		responses, err := b.botService.ListLinks(context.Background(), userID)
		if err != nil {
			return c.Send("❌ Не удалось получить отслеживаемые ссылки")
		}

		if len(responses.Links) == 0 {
			c.Send("Вы не отслеживаете ни одной ссылки")
			return nil
		}

		for _, link := range responses.Links {
			c.Send(link)
		}
		return nil
	}
}