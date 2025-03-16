package grpc

import (
	"context"
	//"tbank/bot/internal/usecase"
	"tbank/bot/api/proto/gen"
	"gopkg.in/telebot.v3"
)


type BotServer struct {
	gen.UnimplementedBotServer
	//usecase 	usecase.UseCase
	telegramBot *telebot.Bot
}

func NewBotServer(bot *telebot.Bot) *BotServer {
	return &BotServer{
		// usecase: usecase,
		telegramBot: bot,
	}
}

func (bs *BotServer) SendUpdate(ctx context.Context, message *gen.CommitUpdate) (*gen.CommitUpdateAnswer, error) {
	return &gen.CommitUpdateAnswer{
		Status: "Succesfully",
	}, nil
}