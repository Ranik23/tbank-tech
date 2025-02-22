package usecase

import (
	"context"
	"tbank/scrapper/config"
	dbmodels "tbank/scrapper/internal/db/models"
	"tbank/scrapper/internal/hub"
	"tbank/scrapper/internal/storage"
)



type UseCase interface {
	RegisterChat(ctx context.Context, chatID int64) 				error
	DeleteChat(ctx context.Context, chatID int64) 					error
	GetLinks(ctx context.Context, chatID int64) 					([]dbmodels.Link, error)
	AddLink(ctx context.Context, link dbmodels.Link, chatID int64) 		error
	RemoveLink(ctx context.Context, link dbmodels.Link, cahtID int64) 	error
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

func (u *UseCaseImpl) RegisterChat(ctx context.Context, chatID int64) error {
	return u.storage.CreateChat(ctx, chatID)
}

func (u *UseCaseImpl) DeleteChat(ctx context.Context, chatID int64) error {

	links, err := u.GetLinks(ctx, chatID)
	if err != nil {
		return err
	}

	for _, link := range links {
		u.hub.RemoveTrack(link, chatID)
	}

	return u.storage.DeleteChat(ctx, chatID)
}

func (u *UseCaseImpl) GetLinks(ctx context.Context, chatID int64) ([]dbmodels.Link, error) {
	return u.storage.GetLinks(ctx, chatID)
}

func (u *UseCaseImpl) AddLink(ctx context.Context, link dbmodels.Link, chatID int64) error {
	u.hub.AddTrack(link, chatID)

	if err := u.storage.CreateLink(ctx, link); err != nil {
		return err
	}

	id, err := u.storage.FindLinkID(ctx, link)
	if err != nil {
		return err
	}

	if err := u.storage.CreateLinkChat(ctx, id, chatID); err != nil {
		return err
	}
	return nil
}

func (u *UseCaseImpl) RemoveLink(ctx context.Context, link dbmodels.Link, chatID int64) error {
	if err := u.hub.RemoveTrack(link, chatID); err != nil {
		return err
	}
	return nil
}
