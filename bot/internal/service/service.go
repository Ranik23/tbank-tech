package service

import (
	"context"
	"fmt"
	"log/slog"
	"tbank/bot/config"
	"tbank/scrapper/api/proto/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Service interface {
	RegisterUser(ctx context.Context, id int64, hashedToken []byte) 							(*gen.RegisterUserResponse, error)
	DeleteUser(ctx context.Context, id int64) 													(*gen.DeleteUserResponse, error) 
	Help(ctx context.Context)																	error 
	AddLink(ctx context.Context, chatID int64, link string, tags []string, filters []string) 	(*gen.LinkResponse, error)
	RemoveLink(ctx context.Context, chatID int64, link string) 									(*gen.LinkResponse, error)
	ListLinks(ctx context.Context, chatID int64) 												(*gen.ListLinksResponse, error) 
}

type service struct {
	config 	*config.Config
	client 	gen.ScrapperClient
	logger 	*slog.Logger
}

func NewService(config *config.Config, logger *slog.Logger) (*service, error) {

	host := config.ScrapperService.Host
	port := config.ScrapperService.Port

	connScrapper, err := grpc.NewClient(fmt.Sprintf("%s:%s", host, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error("Failed to connect to scrapper service", slog.String("error", err.Error()))
		return nil, err
	}

	client := gen.NewScrapperClient(connScrapper)

	return &service{
		config:  config,
		logger:  logger,
		client:  client,
	}, nil
}

func (uc *service) RegisterUser(ctx context.Context, chatID int64, hashedToken []byte) (*gen.RegisterUserResponse, error) {
	uc.logger.Info("Registering chat", slog.Int64("chatID", chatID))

	resp, err := uc.client.RegisterUser(ctx, &gen.RegisterUserRequest{TgUserId: chatID})
	if err != nil {
		uc.logger.Error("failed to register chat", slog.Int64("chatID", chatID), slog.String("error", err.Error()))
		return nil, err
	}

	uc.logger.Info("chat registered successfully", slog.Int64("chatID", chatID))
	return resp, nil
}

func (uc *service) DeleteUser(ctx context.Context, chatID int64) (*gen.DeleteUserResponse, error) {
	uc.logger.Info("Deleting chat", slog.Int64("chatID", chatID))

	resp, err := uc.client.DeleteUser(ctx, &gen.DeleteUserRequest{TgUserId: chatID})
	if err != nil {
		uc.logger.Error("failed to delete chat", slog.Int64("chatID", chatID), slog.String("error", err.Error()))
		return nil, err
	}

	uc.logger.Info("chat deleted successfully", slog.Int64("chatID", chatID))
	return resp, nil
}

func (uc *service) Help(ctx context.Context) error {
	uc.logger.Info("Help requested")
	return nil
}

func (uc *service) AddLink(ctx context.Context, chatID int64, link string, tags []string, filters []string) (*gen.LinkResponse, error) {
	uc.logger.Info("Adding link", slog.Int64("chatID", chatID), slog.String("link", link))

	resp, err := uc.client.AddLink(ctx, &gen.AddLinkRequest{
		TgUserId: chatID,
		Url:     link,
	})
	if err != nil {
		uc.logger.Error("failed to add link", slog.Int64("chatID", chatID), slog.String("link", link), slog.String("error", err.Error()))
		return nil, err
	}

	uc.logger.Info("link added successfully", slog.Int64("chatID", chatID), slog.String("link", link))
	return resp, nil
}

func (uc *service) RemoveLink(ctx context.Context, chatID int64, link string) (*gen.LinkResponse, error) {
	uc.logger.Info("Removing link", slog.Int64("chatID", chatID), slog.String("link", link))

	resp, err := uc.client.RemoveLink(ctx, &gen.RemoveLinkRequest{
		TgUserId: chatID,
		Url:     link,
	})
	if err != nil {
		uc.logger.Error("failed to remove link", slog.Int64("chatID", chatID), slog.String("link", link), slog.String("error", err.Error()))
		return nil, err
	}

	uc.logger.Info("link removed successfully", slog.Int64("chatID", chatID), slog.String("link", link))
	return resp, nil
}

func (uc *service) ListLinks(ctx context.Context, chatID int64) (*gen.ListLinksResponse, error) {
	uc.logger.Info("Listing links", slog.Int64("chatID", chatID))

	resp, err := uc.client.GetLinks(ctx, &gen.GetLinksRequest{
		TgUserId: chatID,
	})
	if err != nil {
		uc.logger.Error("failed to list links", slog.Int64("chatID", chatID), slog.String("error", err.Error()))
		return nil, err
	}

	uc.logger.Info("links listed successfully", slog.Int64("chatID", chatID))
	return resp, nil
}
