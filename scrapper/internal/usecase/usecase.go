package usecase

import (
	"context"
	"tbank/scrapper/config"
	dbmodels "tbank/scrapper/internal/db/models"
	"tbank/scrapper/internal/hub"
	"tbank/scrapper/internal/storage"
)



type UseCase interface {
	RegisterChat(ctx context.Context, chatID uint) 										error
	DeleteChat(ctx context.Context, chatID uint) 										error
	GetLinks(ctx context.Context, chatID uint) 											([]dbmodels.Link, error)
	AddLink(ctx context.Context, link dbmodels.Link, tags []string, filters []string, chatID int64) 	(*dbmodels.Link, error)
	RemoveLink(ctx context.Context, linkID uint) 										error
}

type UseCaseImpl struct {
	cfg 		*config.Config
	storage 	storage.Storage
	hub			*hub.Hub
}

func NewUseCaseImpl(cfg *config.Config, storage storage.Storage, hub *hub.Hub) (*UseCaseImpl, error) {
	return &UseCaseImpl{
		cfg: cfg,
		storage: storage,
		hub: hub,
	}, nil
}

func (u *UseCaseImpl) RegisterChat(ctx context.Context, chatID uint) error {
	return u.storage.CreateChat(ctx, chatID)
}

func (u *UseCaseImpl) DeleteChat(ctx context.Context, chatID uint) error {
	return u.storage.DeleteChat(ctx, chatID)
}

func (u *UseCaseImpl) GetLinks(ctx context.Context, chatID uint) ([]dbmodels.Link, error) {
	return u.storage.GetURLS(ctx, chatID)
}

func (u *UseCaseImpl) AddLink(ctx context.Context, link dbmodels.Link, tags []string, filters []string, chatID int64) (*dbmodels.Link, error) {
	return nil, nil
}

func (u *UseCaseImpl) RemoveLink(ctx context.Context, linkID uint) error {
	return u.storage.DeleteLink(ctx, linkID)
}
