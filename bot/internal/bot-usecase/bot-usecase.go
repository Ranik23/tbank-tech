package botusecase

import (
	"context"
	"fmt"
	"log/slog"
	"tbank/bot/config"
	"tbank/bot/internal/storage"
	"tbank/scrapper/api/proto/gen"
	"google.golang.org/grpc"
)

type UseCase interface {
	RegisterChat(сtx context.Context, id int64) 												(*gen.RegisterChatResponse, error)
	DeleteChat(ctx context.Context, id int64) 													(*gen.DeleteChatResponse, error) 
	Help(ctx context.Context)																	error 
	AddLink(ctx context.Context, chatID int64, link string, tags []string, filters []string) 	(*gen.LinkResponse, error)
	RemoveLink(ctx context.Context, chatID int64, link string) 									(*gen.LinkResponse, error)
	ListLinks(ctx context.Context, chatID int64) 												(*gen.ListLinksResponse, error) 
}


type UseCaseImpl struct {
	config 		*config.Config
	client 		gen.ScrapperClient
	logger 		*slog.Logger
	storage 	storage.Storage
}

func NewUseCaseImpl(config *config.Config, storage storage.Storage, logger *slog.Logger) (*UseCaseImpl, error) {

	host := config.ScrapperService.Host
	port := config.ScrapperService.Port

	connScrapper, err := grpc.NewClient(fmt.Sprintf("%s:%s", host, port),
										grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := gen.NewScrapperClient(connScrapper)

	return &UseCaseImpl{
		config: config,
		storage: storage,
		logger: logger,
		client: client,
	}, nil
}


func (uc *UseCaseImpl) RegisterChat(ctx context.Context, chatID int64) (*gen.RegisterChatResponse, error) {
	return uc.client.RegisterChat(ctx, &gen.RegisterChatRequest{Id: chatID})
}


func (uc *UseCaseImpl) DeleteChat(ctx context.Context, chatID int64) (*gen.DeleteChatResponse, error) {
	return uc.client.DeleteChat(ctx, &gen.DeleteChatRequest{Id: chatID})
}


func (uc *UseCaseImpl) Help(ctx context.Context) error {
	return nil
}

func (uc *UseCaseImpl) AddLink(ctx context.Context, chatID int64, link string, tags []string, filters []string) (*gen.LinkResponse, error) {
	return uc.client.AddLink(ctx, &gen.AddLinkRequest{
		TgChatId: chatID,
		Link: link,
		Tags: tags,
		Filters: filters,
	})
}

func (uc *UseCaseImpl) RemoveLink(ctx context.Context, chatID int64, link string) (*gen.LinkResponse, error) {
	return uc.client.RemoveLink(ctx, &gen.RemoveLinkRequest{
		TgChatId: chatID,
		Link: link,
	})
}

func (uc *UseCaseImpl) ListLinks(ctx context.Context, chatID int64) (*gen.ListLinksResponse, error) {
	return uc.client.GetLinks(ctx, &gen.GetLinksRequest{
		TgChatId: chatID,
	})
}