package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Ranik23/tbank-tech/scrapper/internal/hub"
	dbmodels "github.com/Ranik23/tbank-tech/scrapper/internal/models"
	"github.com/Ranik23/tbank-tech/scrapper/internal/repository"
	"github.com/Ranik23/tbank-tech/scrapper/internal/repository/postgres"
	"github.com/Ranik23/tbank-tech/scrapper/pkg/github_client/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrInvalidLink       = fmt.Errorf("invalid link")
	ErrUserAlreadyExists = fmt.Errorf("user already exists")
	ErrUserNotFound      = fmt.Errorf("user not found")
	ErrLinkNotFound      = fmt.Errorf("link not found")
)

type Service interface {
	RegisterUser(ctx context.Context, userID uint, name string, token string) error
	DeleteUser(ctx context.Context, userID uint) error
	GetLinks(ctx context.Context, userID uint) ([]Link, error)
	AddLink(ctx context.Context, link string, userID uint) error
	RemoveLink(ctx context.Context, link string, userID uint) error
}

type service struct {
	txManager repository.TxManager
	logger    *slog.Logger
	hub       hub.Hub
	repo      repository.Repository
}

func NewService(repo repository.Repository, txManager repository.TxManager,
	hub hub.Hub, logger *slog.Logger) (Service, error) {
	return &service{
		repo:      repo,
		txManager: txManager,
		hub:       hub,
		logger:    logger,
	}, nil
}

func (s *service) RegisterUser(ctx context.Context, userID uint, name string, token string) error {
	var err error
	maxRetries := 3

	for i := 0; i < maxRetries; i++ {
		err = s.txManager.WithTx(ctx, func(txCtx context.Context) error {
			_, err := s.repo.GetUserByID(txCtx, userID)
			if err != nil && !errors.Is(err, postgres.ErrNoUserFound) {
				return err
			}
			if errors.Is(err, postgres.ErrNoUserFound) {
				if err := s.repo.CreateUser(txCtx, userID, name, token); err != nil {
					return err
				}
				return nil
			}
			return ErrUserAlreadyExists
		}, pgx.ReadWrite)

		if err == nil {
			return nil
		}

		if !IsSerializationError(err) {
			s.logger.Error("Non-retryable error", slog.String("error", err.Error()))
			return err
		}

		s.logger.Info("Serialization failure, retrying...", slog.Int("attempt", i+1))
	}

	s.logger.Error("Max retries reached, transaction failed", slog.String("error", err.Error()))
	return fmt.Errorf("max retries reached, transaction failed: %w", err)
}

func IsSerializationError(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.SQLState() == "40001" {
		return true
	}
	return false
}

func (s *service) DeleteUser(ctx context.Context, userID uint) error {
	var err error

	maxRetries := 3

	for i := 1; i < maxRetries; i++ {
		err = s.txManager.WithTx(ctx, func(txCtx context.Context) error {
			_, err := s.repo.GetUserByID(txCtx, userID)
			if err != nil {
				if errors.Is(err, postgres.ErrNoUserFound) {
					return ErrUserNotFound
				}
				return err
			}
			if err := s.repo.DeleteUser(txCtx, userID); err != nil {
				return err
			}
			return nil
		}, pgx.ReadWrite)

		if err == nil {
			return nil
		}

		if !IsSerializationError(err) {
			s.logger.Error("Non-retryable error", slog.String("error", err.Error()))
			return err
		}
		s.logger.Info("Serialization failure, retrying...", slog.Int("attempt", i+1))
	}

	s.logger.Error("Max retries reached, transaction failed", slog.String("error", err.Error()))
	return fmt.Errorf("max retries reached, transaction failed: %w", err)
}

func (s *service) GetLinks(ctx context.Context, userID uint) ([]Link, error) {
	var links []Link
	var err error

	maxRetries := 3

	for i := 1; i <= maxRetries; i++ {
		err = s.txManager.WithTx(ctx, func(txCtx context.Context) error {
			_, err := s.repo.GetUserByID(txCtx, userID)
			if err != nil {
				if errors.Is(err, postgres.ErrNoUserFound) {
					return ErrUserNotFound
				}
				return err
			}

			dbLinks, err := s.repo.GetLinks(txCtx, userID)
			if err != nil {
				return err
			}

			for _, link := range dbLinks {
				links = append(links, Link{
					URL: link.Url,
					ID:  link.ID,
				})
			}

			return nil
		}, pgx.ReadOnly)

		if err == nil {
			return links, nil
		}

		if !IsSerializationError(err) {
			s.logger.Error("Non-retryable error", slog.String("error", err.Error()))
			return nil, err
		}
		s.logger.Info("Serialization failure, retrying...", slog.Int("attempt", i+1))
	}
	
	s.logger.Error("Max retries reached, transaction failed", slog.String("error", err.Error()))
	return nil, fmt.Errorf("max retries reached, transaction failed: %w", err)
}

func (s *service) AddLink(ctx context.Context, link string, userID uint) error {
	_, _, err := utils.GetLinkParams(link)
	if err != nil {
		s.logger.Error("Wrong URL schema", slog.String("error", err.Error()))
		return ErrInvalidLink
	}
	err = s.txManager.WithTx(ctx, func(txCtx context.Context) error {
		user, err := s.repo.GetUserByID(txCtx, userID)
		if err != nil {
			if errors.Is(err, postgres.ErrNoUserFound) {
				return ErrUserNotFound
			}
			return err
		}

		var linkObj *dbmodels.Link

		linkObj, err = s.repo.GetLinkByURL(txCtx, link)
		if err != nil && !errors.Is(err, postgres.ErrNoLinkFound) {
			return err
		}

		if errors.Is(err, postgres.ErrNoLinkFound) {
			if err = s.repo.CreateLink(txCtx, link); err != nil {
				return err
			}
			linkObj, err = s.repo.GetLinkByURL(txCtx, link)
			if err != nil {
				return err
			}
		}

		if err := s.repo.CreateLinkUser(txCtx, linkObj.ID, userID); err != nil {
			return err
		}

		if err := s.hub.AddLink(link, userID, user.Token, 10*time.Second); err != nil {
			return err
		}
		return nil
	}, pgx.ReadWrite)

	if err != nil {
		s.logger.Error("Failed to complete transaction", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func (s *service) RemoveLink(ctx context.Context, link string, userID uint) error {
	_, _, err := utils.GetLinkParams(link)
	if err != nil {
		s.logger.Error("Wrong URL schema", slog.String("link", link), slog.String("error", err.Error()))
		return ErrInvalidLink
	}
	err = s.txManager.WithTx(ctx, func(txCtx context.Context) error {
		linkObj, err := s.repo.GetLinkByURL(txCtx, link)
		if err != nil {
			if errors.Is(err, postgres.ErrNoLinkFound) {
				return ErrLinkNotFound
			}
			return err
		}
		_, err = s.repo.GetUserByID(txCtx, userID)
		if err != nil {
			return ErrUserNotFound
		}

		if err = s.repo.DeleteLink(txCtx, linkObj.ID); err != nil {
			return err
		}

		if err = s.hub.RemoveLink(link, userID); err != nil {
			return err
		}

		return nil
	}, pgx.ReadWrite)
	if err != nil {
		s.logger.Error("Failed to complete transaction", slog.String("error", err.Error()))
		return err
	}
	return nil
}
