package grpc

import (
	"context"
	"github.com/Ranik23/tbank-tech/bot/internal/service"
	genBot"github.com/Ranik23/tbank-tech/bot/api/proto/gen"
	"gopkg.in/telebot.v3"
)

type BotGRPCServer struct {
	genBot.UnimplementedBotServer
	botService  service.Service
	telegramBot *telebot.Bot
}

func NewBotGRPCServer(bot *telebot.Bot, botService service.Service) *BotGRPCServer {
	return &BotGRPCServer{
		botService:  botService,
		telegramBot: bot,
	}
}

func (bs *BotGRPCServer) SendUpdate(ctx context.Context, message *genBot.CommitUpdate) (*genBot.CommitUpdateAnswer, error) {
	return &genBot.CommitUpdateAnswer{
		Status: "Successfully",
	}, nil
}