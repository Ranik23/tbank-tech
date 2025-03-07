package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"tbank/scrapper/config"
	"tbank/scrapper/internal/hub"
	dbmodels "tbank/scrapper/internal/models"
	"tbank/scrapper/internal/storage"
)

//TODO


var (
	ErrEmptyURL = fmt.Errorf("empty url")
)


type UseCase interface {
	RegisterUser(ctx context.Context, userID uint, name string) 										error
	DeleteUser(ctx context.Context, userID uint) 														error
	GetLinks(ctx context.Context, userID uint) 															([]dbmodels.Link, error)
	AddLink(ctx context.Context, link dbmodels.Link, userID uint) 										(*dbmodels.Link, error)
	RemoveLink(ctx context.Context, link dbmodels.Link, userID uint) 									error
}

type UseCaseImpl struct {
	logger 		*slog.Logger
	hub 		*hub.Hub
	cfg 		*config.Config
	storage 	storage.Storage
}

func NewUseCaseImpl(cfg *config.Config, storage storage.Storage,
					hub *hub.Hub, logger *slog.Logger) (*UseCaseImpl, error) {
	return &UseCaseImpl{
		cfg: cfg,
		storage: storage,
		hub: hub,
		logger: logger,
	}, nil
}

func (usecase *UseCaseImpl) RegisterUser(ctx context.Context, userID uint, name string) error {
	return usecase.storage.CreateUser(ctx, userID, name)
}

func (usecase *UseCaseImpl) DeleteUser(ctx context.Context, userID uint) error {
	return usecase.storage.DeleteUser(ctx, userID)
}

func (usecase *UseCaseImpl) GetLinks(ctx context.Context, userID uint) ([]dbmodels.Link, error) {
	return usecase.storage.GetURLS(ctx, userID)
}

func (usecase *UseCaseImpl) AddLink(ctx context.Context, link dbmodels.Link, userID uint) (*dbmodels.Link, error) {

	if link.Url == "" {
		return nil, ErrEmptyURL
	}

	usecase.hub.AddTrack(link.Url, userID)

	if err := usecase.storage.CreateLink(ctx, link.Url); err != nil {
		return nil, err
	}
	
	linkNew, err := usecase.storage.GetLinkByURL(ctx, link.Url)
	if err != nil {
		return nil, err
	}

	if err := usecase.storage.CreateLinkUser(ctx, linkNew.ID, userID); err != nil {
		return nil, err
	}

	return linkNew, nil
}

func (usecase*UseCaseImpl) RemoveLink(ctx context.Context, link dbmodels.Link, userID uint) error {

	if link.Url == "" {
		return ErrEmptyURL
	}

	usecase.hub.RemoveTrack(link.Url, userID)

	linkNew, err := usecase.storage.GetLinkByURL(ctx, link.Url)
	if err != nil {
		return err
	}

	if err := usecase.storage.DeleteLink(ctx, linkNew.ID); err != nil {
		return err
	}
	
	return nil
}
