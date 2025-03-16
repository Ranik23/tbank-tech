package bothandlers

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"tbank/bot/internal/service"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/telebot.v3"
)


func MessageHandler(usecase service.Service, users *sync.Map) telebot.HandlerFunc {
    return func(c telebot.Context) error {
        userID := c.Sender().ID
        chatID := c.Chat().ID

        userRaw, exists := users.Load(userID)
        var user *User
        if !exists {
            user = &User{state: StateFinished}
            users.Store(userID, user)
        } else {
            user = userRaw.(*User)
        }

        text := c.Text()

        switch user.state {
		case StateWaitingForLinkUNLINK:
			user.link = text
			user.state = StateFinished
			users.Store(userID, user)

			response, err := usecase.RemoveLink(context.Background(), chatID, user.link)
			if err != nil {
				return c.Send(fmt.Sprintf("Ошибка: %v", err))
			}
			return c.Send(fmt.Sprintf("Ссылка `%s` успешно отозвана!", response.Url))
        case StateWaitingForLinkLINK:
            user.link = text
            user.state = StateWaitingForTagsLINK
            users.Store(userID, user)
            return c.Send("Введите теги через запятую (например: новости, финансы):")
        case StateWaitingForTagsLINK:
            user.tags = parseList(text)
            user.state = StateWaitingForFiltersLINK
            users.Store(userID, user)
            return c.Send("Введите фильтры через запятую (например: акции, скидки):")
        case StateWaitingForFiltersLINK:
            user.filters = parseList(text)

            response, err := usecase.AddLink(context.Background(), chatID, user.link, user.tags, user.filters)
            user.state = StateFinished
            users.Store(userID, user)

            if err != nil {
                return c.Send(fmt.Sprintf("Ошибка: %v", err))
            }
            return c.Send(fmt.Sprintf("Ссылка '%s' успешно добавлена!", response.Url))
        case StateWaitingForTheToken:

            hashedToken, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
            if err != nil {
                return c.Send(fmt.Sprintf("Ошибка %v", err))
            }

            response, err := usecase.RegisterUser(context.Background(), userID, hashedToken)
            if err != nil {
                return c.Send(fmt.Sprintf("Ошибка %v", err))
            }

            return c.Send(response.GetMessage())
        }

        return nil
    }
}

func parseList(input string) []string {
	parts := strings.Split(input, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}
