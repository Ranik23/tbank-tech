package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	dbmodels "tbank/scrapper/internal/models"
	"tbank/scrapper/internal/repository"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrUniqueViolation = "23505"
)

type postgresRepository struct {
	logger    *slog.Logger
	txManager repository.TxManager
}

func NewPostgresRepository(txManager repository.TxManager, logger *slog.Logger) repository.Repository {
	return &postgresRepository{txManager: txManager, logger: logger}
}

func (s *postgresRepository) GetLinkByURL(ctx context.Context, url string) (*dbmodels.Link, error) {
	s.logger.Info("GetLinkByURL called", slog.String("url", url))

	executor := s.txManager.GetExecutor(ctx)
	var link dbmodels.Link
	var id uint
	query := `SELECT id, url FROM links WHERE url = $1`

	if err := executor.QueryRow(ctx, query, url).Scan(&id, &link.Url); err != nil {
		s.logger.Error("GetLinkByURL failed", slog.String("error", err.Error()))
		return nil, err
	}

	link.ID = id

	s.logger.Info("GetLinkByURL success", slog.Any("link", link))
	return &link, nil
}


func (s *postgresRepository) GetURLS(ctx context.Context, userID uint) ([]dbmodels.Link, error) {
	s.logger.Info("GetURLS called", slog.Uint64("userID", uint64(userID)))

	var links []dbmodels.Link
	executor := s.txManager.GetExecutor(ctx)
	query := `
		SELECT l.id, l.url
		FROM links l
		JOIN link_users lu ON l.id = lu.link_id
		WHERE lu.user_id = $1
	`

	rows, err := executor.Query(ctx, query, userID)
	if err != nil {
		s.logger.Error("GetURLS failed", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to query links: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var link dbmodels.Link
		if err := rows.Scan(&link.ID, &link.Url); err != nil {
			s.logger.Error("GetURLS scan failed", slog.String("error", err.Error()))
			return nil, fmt.Errorf("failed to scan link: %w", err)
		}
		links = append(links, link)
	}

	if err := rows.Err(); err != nil {
		s.logger.Error("GetURLS rows error", slog.String("error", err.Error()))
		return nil, fmt.Errorf("rows error: %w", err)
	}

	s.logger.Info("GetURLS success", slog.Int("count", len(links)))
	return links, nil
}

func (s *postgresRepository) DeleteUser(ctx context.Context, userID uint) error {
	s.logger.Info("DeleteUser called", slog.Uint64("userID", uint64(userID)))

	executor := s.txManager.GetExecutor(ctx)
	query := `DELETE FROM users WHERE user_id = $1`

	_, err := executor.Exec(ctx, query, userID)
	if err != nil {
		s.logger.Error("DeleteUser failed", slog.String("error", err.Error()))
		return fmt.Errorf("failed to delete user: %w", err)
	}

	s.logger.Info("DeleteUser success")
	return nil
}

func (s *postgresRepository) DeleteLink(ctx context.Context, linkID uint) error {
	s.logger.Info("DeleteLink called", slog.Uint64("linkID", uint64(linkID)))

	executor := s.txManager.GetExecutor(ctx)
	query := `DELETE FROM links WHERE id = $1`

	_, err := executor.Exec(ctx, query, linkID)
	if err != nil {
		s.logger.Error("DeleteLink failed", slog.String("error", err.Error()))
		return fmt.Errorf("failed to delete link: %w", err)
	}

	s.logger.Info("DeleteLink success")
	return nil
}

func (s *postgresRepository) DeleteLinkUser(ctx context.Context, linkID uint, userID uint) error {
	s.logger.Info("DeleteLinkUser called", slog.Uint64("linkID", uint64(linkID)), slog.Uint64("userID", uint64(userID)))

	executor := s.txManager.GetExecutor(ctx)
	query := `DELETE FROM link_users WHERE link_id = $1 AND user_id = $2`

	_, err := executor.Exec(ctx, query, linkID, userID)
	if err != nil {
		s.logger.Error("DeleteLinkUser failed", slog.String("error", err.Error()))
		return fmt.Errorf("failed to delete link-user association: %w", err)
	}

	s.logger.Info("DeleteLinkUser success")
	return nil
}

func (s *postgresRepository) CreateLinkUser(ctx context.Context, linkID uint, userID uint) error {
	s.logger.Info("CreateLinkUser called", slog.Uint64("linkID", uint64(linkID)), slog.Uint64("userID", uint64(userID)))

	executor := s.txManager.GetExecutor(ctx)
	query := `INSERT INTO link_users (link_id, user_id) VALUES ($1, $2)`

	_, err := executor.Exec(ctx, query, linkID, userID)
	if err != nil {
		s.logger.Error("CreateLinkUser failed", slog.String("error", err.Error()))
		return fmt.Errorf("failed to create link-user association: %w", err)
	}

	s.logger.Info("CreateLinkUser success")
	return nil
}

func (s *postgresRepository) CreateUser(ctx context.Context, userID uint, name string) error {
	s.logger.Info("CreateUser called", slog.Uint64("userID", uint64(userID)), slog.String("name", name))

	executor := s.txManager.GetExecutor(ctx)
	query := `INSERT INTO users (user_id, name) VALUES ($1, $2)`

	_, err := executor.Exec(ctx, query, userID, name)
	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) && pgError.Code == ErrUniqueViolation {
			s.logger.Warn("CreateUser duplicate, skipping")
			return nil
		}
		s.logger.Error("CreateUser failed", slog.String("error", err.Error()))
		return fmt.Errorf("failed to create user: %w", err)
	}

	s.logger.Info("CreateUser success")
	return nil
}

func (s *postgresRepository) CreateLink(ctx context.Context, link string) error {
	s.logger.Info("CreateLink called", slog.String("link", link))

	executor := s.txManager.GetExecutor(ctx)
	query := `INSERT INTO links (url) VALUES ($1)`

	_, err := executor.Exec(ctx, query, link)
	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) && pgError.Code == ErrUniqueViolation {
			s.logger.Warn("CreateLink duplicate, skipping")
			return nil
		}
		s.logger.Error("CreateLink failed", slog.String("error", err.Error()))
		return fmt.Errorf("failed to create link: %w", err)
	}

	s.logger.Info("CreateLink success")
	return nil
}

func (s *postgresRepository) GetLinkByID(ctx context.Context, id uint) (*dbmodels.Link, error) {
	s.logger.Info("GetLinkByID called", slog.Uint64("id", uint64(id)))

	executor := s.txManager.GetExecutor(ctx)
	var link dbmodels.Link
	query := `SELECT id, url FROM links WHERE id = $1`

	err := executor.QueryRow(ctx, query, id).Scan(&link.ID, &link.Url)
	if err != nil {
		s.logger.Error("GetLinkByID failed", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to get link by ID: %w", err)
	}

	s.logger.Info("GetLinkByID success", slog.Any("link", link))
	return &link, nil
}
