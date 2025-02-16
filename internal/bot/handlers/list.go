package handlers

import (
	"tbank/internal/usecase"
	"gopkg.in/telebot.v3"
)

// ListHandler показывает список отслеживаемых ссылок
type ListHandler struct {
	usecase usecase.UseCase
}

func NewListHandler(usecase usecase.UseCase) *ListHandler {
	return &ListHandler{usecase: usecase}
}

func (h *ListHandler) Handle(c telebot.Context) error {
	userID := (int)(c.Sender().ID)
	
	h.usecase.ListLinks(userID)
	return c.Send("Вот ваш список отслеживаемых ссылок.")
}
