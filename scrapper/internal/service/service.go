package service

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
	EmptyLinkURL = ""
)

type Service interface {
	RegisterUser(ctx context.Context, userID uint, name string) 										error
	DeleteUser(ctx context.Context, userID uint) 														error
	GetLinks(ctx context.Context, userID uint) 															([]dbmodels.Link, error)
	AddLink(ctx context.Context, link dbmodels.Link, userID uint) 										(*dbmodels.Link, error)
	RemoveLink(ctx context.Context, link dbmodels.Link, userID uint) 									error
}

type service struct {
	logger 		*slog.Logger
	hub			hub.Hub
	repo		repository.Repository
}

func NewService(repo repository.Repository, hub hub.Hub, logger *slog.Logger) (Service, error) {
	return &service{
		repo: repo,
		hub: hub,
		logger: logger,
	}, nil
}

func (s *service) RegisterUser(ctx context.Context, userID uint, name string) error {
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				s.logger.Error("transaction rollback failed", "error", rollbackErr)
				err = fmt.Errorf("original error: %w, rollback error: %w", err, rollbackErr)
			}
		}
	}()

	if err = s.repo.CreateUser(ctx, userID, name); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *service) DeleteUser(ctx context.Context, userID uint) error {
	return s.repo.DeleteUser(ctx, userID)
}

func (s *service) GetLinks(ctx context.Context, userID uint) ([]dbmodels.Link, error) {
	return s.repo.GetURLS(ctx, userID)
}

func (s *service) AddLink(ctx context.Context, link dbmodels.Link, userID uint) (*dbmodels.Link, error) {
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				s.logger.Error("transaction rollback failed", "error", rollbackErr)
				err = fmt.Errorf("original error: %w, rollback error: %w", err, rollbackErr)
			}
		}
	}()

	if link.Url == EmptyLinkURL {
		err = ErrEmptyLink
		return nil, err
	}

	s.hub.AddLink(link.Url, userID)

	if err = s.repo.CreateLink(ctx, link.Url); err != nil {
		return nil, err
	}

	linkNew, err := s.repo.GetLinkByURL(ctx, link.Url)
	if err != nil {
		return nil, err
	}

	if err = s.repo.CreateLinkUser(ctx, linkNew.ID, userID); err != nil {
		return nil, err
	}
	
	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return linkNew, nil
}


func (s *service) RemoveLink(ctx context.Context, link dbmodels.Link, userID uint) error {
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				s.logger.Error("transaction rollback failed", "error", rollbackErr)
				err = fmt.Errorf("original error: %w, rollback error: %w", err, rollbackErr)
			}
		}
	}()

	if link.Url == "" {
		err = ErrEmptyLink
		return err
	}

	s.hub.RemoveLink(link.Url, userID)

	linkNew, err := s.repo.GetLinkByURL(ctx, link.Url)
	if err != nil {
		return err
	}

	if err = s.repo.DeleteLink(ctx, linkNew.ID); err != nil {
		return err
	}
	
	return tx.Commit(ctx)
}
