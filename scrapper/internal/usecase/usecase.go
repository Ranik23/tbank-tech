package usecase

import (
	"context"
	"fmt"
	"tbank/bot/api/proto/gen"
	"tbank/scrapper/config"
	dbmodels "tbank/scrapper/internal/db/models"
	"tbank/scrapper/internal/storage"
	"time"
	gocron "github.com/go-co-op/gocron/v2"
	"google.golang.org/grpc"
)



type UseCase interface {
	RegisterChat(ctx context.Context, chatID uint) 										error
	DeleteChat(ctx context.Context, chatID uint) 										error
	GetLinks(ctx context.Context, chatID uint) 											([]dbmodels.Link, error)
	AddLink(ctx context.Context, link dbmodels.Link, tags []string, filters []string) 	(*dbmodels.Link, error)
	RemoveLink(ctx context.Context, linkID uint) 										error
}

type UseCaseImpl struct {
	cfg 		*config.Config
	storage 	storage.Storage
	scheduler 	gocron.Scheduler
	client 		gen.BotClient
}

func NewUseCaseImpl(cfg *config.Config, storage storage.Storage, scheduler gocron.Scheduler) (*UseCaseImpl, error) {

	host := cfg.Bot.Host
	port := cfg.Bot.Port

	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", host, port), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := gen.NewBotClient(conn)

	return &UseCaseImpl{
		cfg: cfg,
		storage: storage,
		scheduler: scheduler,
		client: client,
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

func (u *UseCaseImpl) AddLink(ctx context.Context, link dbmodels.Link, tags []string, filters []string) (*dbmodels.Link, error) {
	_, err := u.scheduler.NewJob(gocron.DurationJob(
								10 * time.Second,
						),
						gocron.NewTask(
							func(client gen.BotClient) {
								client.SendUpdate(ctx, &gen.UpdateMessage{})
							},
							u.client,
						),
					)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (u *UseCaseImpl) RemoveLink(ctx context.Context, linkID uint) error {
	return u.storage.DeleteLink(ctx, linkID)
}
