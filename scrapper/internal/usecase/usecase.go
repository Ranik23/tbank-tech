package usecase

import (
	"context"
	"fmt"
	"log/slog"
	dbmodels "tbank/scrapper/internal/models"
	"tbank/scrapper/internal/hub"
	"tbank/scrapper/internal/repository"
)


var (
	ErrEmptyLink = fmt.Errorf("empty link")
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
	hub			hub.Hub
	repo		repository.Repository
}

func NewUseCaseImpl(repo repository.Repository, hub hub.Hub, logger *slog.Logger) (UseCase, error) {
	return &usecaseImpl{
		repo: repo,
		hub: hub,
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
	tx, err := usecase.repo.BeginTx(ctx) // может все-таки явно указать что это менеджер?
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

	usecase.hub.AddLink(link.Url, userID)

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

	usecase.hub.RemoveLink(link.Url, userID)

	linkNew, err := usecase.repo.GetLinkByURL(ctx, link.Url)
	if err != nil {
		return err
	}

	if err := usecase.repo.DeleteLink(ctx, linkNew.ID); err != nil {
		return err
	}
	
	return nil
}
