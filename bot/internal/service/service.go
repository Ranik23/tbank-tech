package service

import (
	"context"
	"log/slog"
	"github.com/Ranik23/tbank-tech/scrapper/api/proto/gen"

	"google.golang.org/grpc"
)

type Service interface {
	RegisterUser(ctx context.Context, id int64, name string, hashedToken string) (*gen.RegisterUserResponse, error)
	DeleteUser(ctx context.Context, id int64) 									 (*gen.DeleteUserResponse, error)
	Help(ctx context.Context) error
	AddLink(ctx context.Context, userID int64, link string) 					 (*gen.LinkResponse, error)
	RemoveLink(ctx context.Context, userID int64, link string) 					 (*gen.LinkResponse, error)
	ListLinks(ctx context.Context, userID int64) 							     (*gen.ListLinksResponse, error)
}

type service struct {
	client gen.ScrapperClient
	logger *slog.Logger
}

func NewService(connScrapper grpc.ClientConnInterface, logger *slog.Logger) Service {
	client := gen.NewScrapperClient(connScrapper)
	return &service{
		logger: logger,
		client: client,
	}
}

func (uc *service) RegisterUser(ctx context.Context, userID int64, name string, hashedToken string) (*gen.RegisterUserResponse, error) {
	uc.logger.Info("Registering user", slog.Int64("userID", userID))

	resp, err := uc.client.RegisterUser(ctx, &gen.RegisterUserRequest{
		TgUserId: userID,
		Name:     name,
		Token:    hashedToken,
	})

	if err != nil {
		uc.logger.Error("Failed to register user", slog.Int64("userID", userID), slog.String("error", err.Error()))
		return nil, err
	}

	uc.logger.Info("User registered successfully", slog.Int64("userID", userID))
	return resp, nil
}

func (uc *service) DeleteUser(ctx context.Context, userID int64) (*gen.DeleteUserResponse, error) {
	uc.logger.Info("Deleting user", slog.Int64("userID", userID))

	resp, err := uc.client.DeleteUser(ctx, &gen.DeleteUserRequest{TgUserId: userID})
	if err != nil {
		uc.logger.Error("Failed to delete chat", slog.Int64("userID", userID), slog.String("error", err.Error()))
		return nil, err
	}

	uc.logger.Info("Chat deleted successfully", slog.Int64("userID", userID))
	return resp, nil
}

func (uc *service) Help(ctx context.Context) error {
	uc.logger.Info("Help requested")
	return nil
}

func (uc *service) AddLink(ctx context.Context, userID int64, link string) (*gen.LinkResponse, error) {
	uc.logger.Info("Adding link", slog.Int64("userID", userID), slog.String("link", link))

	resp, err := uc.client.AddLink(ctx, &gen.AddLinkRequest{
		TgUserId: userID,
		Url:      link,
	})
	if err != nil {
		uc.logger.Error("Failed to add link", slog.Int64("userID", userID), slog.String("link", link), slog.String("error", err.Error()))
		return nil, err
	}

	uc.logger.Info("Link added successfully", slog.Int64("userID", userID), slog.String("link", link))
	return resp, nil
}

func (uc *service) RemoveLink(ctx context.Context, userID int64, link string) (*gen.LinkResponse, error) {
	uc.logger.Info("Removing link", slog.Int64("userID", userID), slog.String("link", link))

	resp, err := uc.client.RemoveLink(ctx, &gen.RemoveLinkRequest{
		TgUserId: userID,
		Url:      link,
	})
	if err != nil {
		uc.logger.Error("Failed to remove link", slog.Int64("userID", userID), slog.String("link", link), slog.String("error", err.Error()))
		return nil, err
	}

	uc.logger.Info("Link removed successfully", slog.Int64("userID", userID), slog.String("link", link))
	return resp, nil
}

func (uc *service) ListLinks(ctx context.Context, userID int64) (*gen.ListLinksResponse, error) {
	uc.logger.Info("Listing links", slog.Int64("userID", userID))

	resp, err := uc.client.GetLinks(ctx, &gen.GetLinksRequest{
		TgUserId: userID,
	})
	if err != nil {
		uc.logger.Error("Failed to list links", slog.Int64("userID", userID), slog.String("error", err.Error()))
		return nil, err
	}

	uc.logger.Info("Links listed successfully", slog.Int64("userID", userID))
	return resp, nil
}
