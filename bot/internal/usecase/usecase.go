package usecase

import (
	"context"

	"gopkg.in/telebot.v3"
)


type UseCase interface {
	SendMessage(ctx context.Context, recipients ...*telebot.Chat) error
}

type UseCaseImpl struct {
	bot *telebot.Bot
}

func NewUseCaseImp(bot *telebot.Bot) *UseCaseImpl {
	return &UseCaseImpl{
		bot: bot,
	}
}

func (uc *UseCaseImpl) SendMessage(ctx context.Context, recipients ...*telebot.Chat) error {
	for _, chat := range recipients {
		uc.bot.Send(chat, "Hello") // тест
	}
	return nil
}