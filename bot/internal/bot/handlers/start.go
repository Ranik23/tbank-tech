package handlers

import (
	"tbank/bot/internal/usecase"
	"gopkg.in/telebot.v3"
)

// StartHandler отвечает за команду /start
type StartHandler struct {
	usecase usecase.UseCase
}

func NewStartHandler(usecase usecase.UseCase) *StartHandler {
	return &StartHandler{usecase: usecase}
}

func (h *StartHandler) Handle(c telebot.Context) error {
	return c.Send("Привет! Я бот для трекинга ссылок.")
}
