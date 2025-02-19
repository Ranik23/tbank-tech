package botusecase

import (
	"context"
	"fmt"
	"tbank/bot/config"
	"tbank/bot/internal/storage"
	"tbank/scrapper/api/proto/gen"
	"google.golang.org/grpc"
	"log/slog"
)

type UseCase interface {
	RegisterChat(ctx context.Context, id int64) 												(*gen.RegisterChatResponse, error)
	DeleteChat(ctx context.Context, id int64) 													(*gen.DeleteChatResponse, error) 
	Help(ctx context.Context)																	error 
	AddLink(ctx context.Context, chatID int64, link string, tags []string, filters []string) 	(*gen.LinkResponse, error)
	RemoveLink(ctx context.Context, chatID int64, link string) 									(*gen.LinkResponse, error)
	ListLinks(ctx context.Context, chatID int64) 												(*gen.ListLinksResponse, error) 
}

type UseCaseImpl struct {
	config 	*config.Config
	client 	gen.ScrapperClient
	logger 	*slog.Logger
	storage storage.Storage
}

func NewUseCaseImpl(config *config.Config, storage storage.Storage, logger *slog.Logger) (*UseCaseImpl, error) {

	host := config.ScrapperService.Host
	port := config.ScrapperService.Port

	connScrapper, err := grpc.NewClient(fmt.Sprintf("%s:%s", host, port), grpc.WithInsecure())
	if err != nil {
		logger.Error("failed to connect to scrapper service", slog.String("error", err.Error()))
		return nil, err
	}

	client := gen.NewScrapperClient(connScrapper)

	logger.Info("scrapper service client created", slog.String("host", host), slog.String("port", port))

	return &UseCaseImpl{
		config:  config,
		storage: storage,
		logger:  logger,
		client:  client,
	}, nil
}

func (uc *UseCaseImpl) RegisterChat(ctx context.Context, chatID int64) (*gen.RegisterChatResponse, error) {
	uc.logger.Info("Registering chat", slog.Int64("chatID", chatID))

	resp, err := uc.client.RegisterChat(ctx, &gen.RegisterChatRequest{Id: chatID})
	if err != nil {
		uc.logger.Error("failed to register chat", slog.Int64("chatID", chatID), slog.String("error", err.Error()))
		return nil, err
	}

	uc.logger.Info("chat registered successfully", slog.Int64("chatID", chatID))
	return resp, nil
}

func (uc *UseCaseImpl) DeleteChat(ctx context.Context, chatID int64) (*gen.DeleteChatResponse, error) {
	uc.logger.Info("Deleting chat", slog.Int64("chatID", chatID))

	resp, err := uc.client.DeleteChat(ctx, &gen.DeleteChatRequest{Id: chatID})
	if err != nil {
		uc.logger.Error("failed to delete chat", slog.Int64("chatID", chatID), slog.String("error", err.Error()))
		return nil, err
	}

	uc.logger.Info("chat deleted successfully", slog.Int64("chatID", chatID))
	return resp, nil
}

func (uc *UseCaseImpl) Help(ctx context.Context) error {
	uc.logger.Info("Help requested")
	return nil
}

func (uc *UseCaseImpl) AddLink(ctx context.Context, chatID int64, link string, tags []string, filters []string) (*gen.LinkResponse, error) {
	uc.logger.Info("Adding link", slog.Int64("chatID", chatID), slog.String("link", link))

	resp, err := uc.client.AddLink(ctx, &gen.AddLinkRequest{
		TgChatId: chatID,
		Link:     link,
		Tags:     tags,
		Filters:  filters,
	})
	if err != nil {
		uc.logger.Error("failed to add link", slog.Int64("chatID", chatID), slog.String("link", link), slog.String("error", err.Error()))
		return nil, err
	}

	uc.logger.Info("link added successfully", slog.Int64("chatID", chatID), slog.String("link", link))
	return resp, nil
}

func (uc *UseCaseImpl) RemoveLink(ctx context.Context, chatID int64, link string) (*gen.LinkResponse, error) {
	uc.logger.Info("Removing link", slog.Int64("chatID", chatID), slog.String("link", link))

	resp, err := uc.client.RemoveLink(ctx, &gen.RemoveLinkRequest{
		TgChatId: chatID,
		Link:     link,
	})
	if err != nil {
		uc.logger.Error("failed to remove link", slog.Int64("chatID", chatID), slog.String("link", link), slog.String("error", err.Error()))
		return nil, err
	}

	uc.logger.Info("link removed successfully", slog.Int64("chatID", chatID), slog.String("link", link))
	return resp, nil
}

func (uc *UseCaseImpl) ListLinks(ctx context.Context, chatID int64) (*gen.ListLinksResponse, error) {
	uc.logger.Info("Listing links", slog.Int64("chatID", chatID))

	resp, err := uc.client.GetLinks(ctx, &gen.GetLinksRequest{
		TgChatId: chatID,
	})
	if err != nil {
		uc.logger.Error("failed to list links", slog.Int64("chatID", chatID), slog.String("error", err.Error()))
		return nil, err
	}

	uc.logger.Info("links listed successfully", slog.Int64("chatID", chatID))
	return resp, nil
}
