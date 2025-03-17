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
	RegisterUser(ctx context.Context, userID uint, name string) error
	DeleteUser(ctx context.Context, userID uint) error
	GetLinks(ctx context.Context, userID uint) ([]dbmodels.Link, error)
	AddLink(ctx context.Context, link dbmodels.Link, userID uint) (*dbmodels.Link, error)
	RemoveLink(ctx context.Context, link dbmodels.Link, userID uint) error
}

type service struct {
	txManager repository.TxManager
	logger   *slog.Logger
	hub      hub.Hub
	repo     repository.Repository
}

func NewService(repo repository.Repository, txManager repository.TxManager, hub hub.Hub, logger *slog.Logger) (Service, error) {
	return &service{
		repo:      repo,
		txManager: txManager,
		hub:       hub,
		logger:    logger,
	}, nil
}

func (s *service) RegisterUser(ctx context.Context, userID uint, name string) error {
	return s.txManager.WithTx(ctx, func(txCtx context.Context) error {
		return s.repo.CreateUser(txCtx, userID, name)
	})
}

func (s *service) DeleteUser(ctx context.Context, userID uint) error {
	return s.txManager.WithTx(ctx, func(txCtx context.Context) error {
		return s.repo.DeleteUser(txCtx, userID)
	})
}

func (s *service) GetLinks(ctx context.Context, userID uint) ([]dbmodels.Link, error) {
	var links []dbmodels.Link
	err := s.txManager.WithTx(ctx, func(txCtx context.Context) error {
		var err error
		links, err = s.repo.GetURLS(txCtx, userID)
		return err
	})
	return links, err
}

func (s *service) AddLink(ctx context.Context, link dbmodels.Link, userID uint) (*dbmodels.Link, error) {
	var linkNew *dbmodels.Link

	err := s.txManager.WithTx(ctx, func(txCtx context.Context) error {
		if link.Url == EmptyLinkURL {
			return ErrEmptyLink
		}

		if err := s.hub.AddLink(link.Url, userID); err != nil {
			return err
		}

		if err := s.repo.CreateLink(txCtx, link.Url); err != nil {
			return err
		}

		if err := s.repo.CreateUser(txCtx, userID, "random"); err != nil {
			return err
		}

		var err error
		linkNew, err = s.repo.GetLinkByURL(txCtx, link.Url)
		if err != nil {
			return err
		}

		return s.repo.CreateLinkUser(txCtx, linkNew.ID, userID)
	})

	return linkNew, err
}

func (s *service) RemoveLink(ctx context.Context, link dbmodels.Link, userID uint) error {
	return s.txManager.WithTx(ctx, func(txCtx context.Context) error {
		if link.Url == EmptyLinkURL {
			return ErrEmptyLink
		}

		if err := s.hub.RemoveLink(link.Url, userID); err != nil {
			return err
		}

		linkNew, err := s.repo.GetLinkByURL(txCtx, link.Url)
		if err != nil {
			return err
		}

		return s.repo.DeleteLink(txCtx, linkNew.ID)
	})
}
