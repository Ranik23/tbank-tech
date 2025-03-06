package usecase

import (
	"context"
	"tbank/scrapper/config"
	dbmodels "tbank/scrapper/internal/models"
	"tbank/scrapper/internal/hub"
	"tbank/scrapper/internal/storage"
)


//TODO

type UseCase interface {
	RegisterUser(ctx context.Context, userID uint, name string) 										error
	DeleteUser(ctx context.Context, userID uint) 														error
	GetLinks(ctx context.Context, userID uint) 															([]dbmodels.Link, error)
	AddLink(ctx context.Context, link dbmodels.Link, userID int64) 										(*dbmodels.Link, error)
	RemoveLink(ctx context.Context, linkID uint) 														error
}

type UseCaseImpl struct {
	hub 		*hub.Hub
	cfg 		*config.Config
	storage 	storage.Storage
}

func NewUseCaseImpl(cfg *config.Config, storage storage.Storage, hub *hub.Hub) (*UseCaseImpl, error) {
	return &UseCaseImpl{
		cfg: cfg,
		storage: storage,
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

func (usecase *UseCaseImpl) AddLink(ctx context.Context, link dbmodels.Link, userID int64) (*dbmodels.Link, error) {
	usecase.hub.AddTrack(link.Url)
	usecase.storage.CreateLink(ctx, link.Url)

	return nil, nil
}

func (u *UseCaseImpl) RemoveLink(ctx context.Context, linkID uint) error {

	var link string

	u.storage.


	u.hub.RemoveTrack()
	return nil
}
