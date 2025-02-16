package handlers

import (
	"tbank/bot/internal/usecase"
	"gopkg.in/telebot.v3"
)

type HelpHandler struct {
	usecase usecase.UseCase
}

func NewHelpHandler(usecase usecase.UseCase) *HelpHandler {
	return &HelpHandler{usecase: usecase}
}

func (h *HelpHandler) Handle(c telebot.Context) error {
	return c.Send("Команды: /start, /track, /untrack, /list")
}
