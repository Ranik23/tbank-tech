package telegram

import (
	"sync"
	"github.com/Ranik23/tbank-tech/bot/internal/service"
)

type BotHandlers struct {
	botService service.Service
	users	   *sync.Map
}

func NewBotHandler(botService service.Service, users *sync.Map) *BotHandlers {
	return &BotHandlers{
		botService: botService,
		users: users,
	}
}
