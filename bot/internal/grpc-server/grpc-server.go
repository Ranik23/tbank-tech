package grpcserver

import (
	"context"
	botusecase"tbank/bot/internal/bot-usecase"
	"tbank/bot/proto/gen"
	"gopkg.in/telebot.v3"
)


type BotServer struct {
	gen.UnimplementedBotServer
	usecase 	botusecase.UseCase
	telegramBot *telebot.Bot
}

func NewBotServer(usecase botusecase.UseCase, bot *telebot.Bot) *BotServer {
	return &BotServer{
		usecase: usecase,
		telegramBot: bot,
	}
}

func (bs *BotServer) SendUpdate(ctx context.Context, message *gen.UpdateMessage) (*gen.Response, error) {
	return &gen.Response{
		Message: "Succesfully",
	}, nil
}