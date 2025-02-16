package handlers

import (
	"tbank/internal/usecase"
	"gopkg.in/telebot.v3"
)

// UntrackHandler убирает отслеживание ссылки
type UntrackHandler struct {
	usecase usecase.UseCase
}

func NewUntrackHandler(usecase usecase.UseCase) *UntrackHandler {
	return &UntrackHandler{usecase: usecase}
}

func (h *UntrackHandler) Handle(c telebot.Context) error {
	return c.Send("Ссылка удалена из отслеживания.")
}
