package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"tbank/scrapper/config"
	dbmodels "tbank/scrapper/internal/models"
	scheduler "tbank/scrapper/internal/hub"
	"tbank/scrapper/internal/storage/db/postgres"
)


var (
	ErrEmptyLink = fmt.Errorf("empty link")
)

const (
	EmptyLink = ""
)


type UseCase interface {
	RegisterUser(ctx context.Context, userID uint, name string) 										error
	DeleteUser(ctx context.Context, userID uint) 														error
	GetLinks(ctx context.Context, userID uint) 															([]dbmodels.Link, error)
	AddLink(ctx context.Context, link dbmodels.Link, userID uint) 										(*dbmodels.Link, error)
	RemoveLink(ctx context.Context, link dbmodels.Link, userID uint) 									error
}

type usecaseImpl struct {
	logger 		*slog.Logger
	scheduler	*scheduler.Hub
	cfg 		*config.Config
	repo		db.Repository
}

func NewUseCaseImpl(cfg *config.Config, repository db.Repository, scheduler *scheduler.Hub, logger *slog.Logger) (UseCase, error) {
	return &usecaseImpl{
		cfg: cfg,
		repo: repository,
		scheduler: scheduler,
		logger: logger,
	}, nil
}

func (usecase *usecaseImpl) RegisterUser(ctx context.Context, userID uint, name string) error {
	tx, err := usecase.repo.BeginTx(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				usecase.logger.Error("transaction rollback failed", "error", rollbackErr)
			}
		}
	}()

	if err := usecase.repo.CreateUser(ctx, userID, name); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (usecase *usecaseImpl) DeleteUser(ctx context.Context, userID uint) error {
	return usecase.repo.DeleteUser(ctx, userID)
}

func (usecase *usecaseImpl) GetLinks(ctx context.Context, userID uint) ([]dbmodels.Link, error) {
	return usecase.repo.GetURLS(ctx, userID)
}

func (usecase *usecaseImpl) AddLink(ctx context.Context, link dbmodels.Link, userID uint) (*dbmodels.Link, error) {
	tx, err := usecase.repo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				usecase.logger.Error("transaction rollback failed", "error", rollbackErr)
			}
		}
	}()

	if link.Url == EmptyLink {
		err := ErrEmptyLink
		return nil, err
	}

	usecase.scheduler.AddLink(link.Url, userID)

	if err := usecase.repo.CreateLink(ctx, link.Url); err != nil {
		return nil, err
	}

	linkNew, err := usecase.repo.GetLinkByURL(ctx, link.Url)
	if err != nil {
		return nil, err
	}

	if err := usecase.repo.CreateLinkUser(ctx, linkNew.ID, userID); err != nil {
		return nil, err
	}
	
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return linkNew, nil
}


func (usecase*usecaseImpl) RemoveLink(ctx context.Context, link dbmodels.Link, userID uint) error {

	tx, err := usecase.repo.BeginTx(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				usecase.logger.Error("transaction rollback failed", "error", rollbackErr)
			}
		}
	}()

	if link.Url == "" {
		err := ErrEmptyLink
		return err
	}

	usecase.scheduler.RemoveLink(link.Url, userID)

	linkNew, err := usecase.repo.GetLinkByURL(ctx, link.Url)
	if err != nil {
		return err
	}

	if err := usecase.repo.DeleteLink(ctx, linkNew.ID); err != nil {
		return err
	}
	
	return nil
}
