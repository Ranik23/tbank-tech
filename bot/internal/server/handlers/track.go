package handlers

import (
	"tbank/bot/internal/usecase"
	"gopkg.in/telebot.v3"
)

// TrackHandler добавляет отслеживание ссылки
type TrackHandler struct {
	usecase usecase.UseCase
}

func NewTrackHandler(usecase usecase.UseCase) *TrackHandler {
	return &TrackHandler{usecase: usecase}
}

func (h *TrackHandler) Handle(c telebot.Context) error {
	return c.Send("Ссылка добавлена в отслеживание.")
}
