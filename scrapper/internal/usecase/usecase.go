package usecase

import (
	"context"
	dbmodels "tbank/scrapper/internal/db/models"
	"tbank/scrapper/internal/storage"
	gocron "github.com/go-co-op/gocron/v2"
)



type UseCase interface {
	RegisterChat(ctx context.Context, chatID uint) error
	DeleteChat(ctx context.Context, chatID uint) error
	GetLinks(ctx context.Context, chatID uint) ([]dbmodels.Link, error)
	AddLink(ctx context.Context, link dbmodels.Link, tags []string, filters []string) (*dbmodels.Link, error)
	RemoveLink(ctx context.Context, linkID uint) error
}

type UseCaseImpl struct {
	storage storage.Storage
	scheduler gocron.Scheduler
}

func NewUseCase(storage storage.Storage, scheduler gocron.Scheduler) *UseCaseImpl {
	return &UseCaseImpl{storage: storage}
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

func (u *UseCaseImpl) AddLink(ctx context.Context, link dbmodels.Link, tags []string, filters []string) (*dbmodels.Link, error) {
	return nil, nil
}

func (u *UseCaseImpl) RemoveLink(ctx context.Context, linkID uint) error {
	return u.storage.DeleteLink(ctx, linkID)
}
