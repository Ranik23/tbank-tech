package handlers

import (
	"tbank/bot/internal/usecase"
	"gopkg.in/telebot.v3"
)

type ListHandler struct {
	usecase usecase.UseCase
}

func NewListHandler(usecase usecase.UseCase) *ListHandler {
	return &ListHandler{usecase: usecase}
}

func (h *ListHandler) Handle(c telebot.Context) error {
	_ = (int)(c.Sender().ID)
	return c.Send("Вот ваш список отслеживаемых ссылок.")
}
