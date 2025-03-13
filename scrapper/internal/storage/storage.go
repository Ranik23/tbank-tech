package storage

import (
	"context"
	"errors"
	"fmt"
	"tbank/scrapper/config"
	dbmodels "tbank/scrapper/internal/models"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrUniqueViolation = "23505"
)


type Storage interface {
	CreateLink(ctx context.Context, link string) 					error
	CreateUser(ctx context.Context, userID uint, name string) 		error
	CreateLinkUser(ctx context.Context, linkID uint, userID uint) 	error

	DeleteUser(ctx context.Context, userID uint) 					error
	DeleteLink(ctx context.Context, linkID uint) 					error
	DeleteLinkUser(ctx context.Context, linkID uint, userID uint) 	error

	GetURLS(ctx context.Context, userID uint) 						([]dbmodels.Link, error)
	GetLinkByID(ctx context.Context, id uint) 						(*dbmodels.Link, error)
	GetLinkByURL(ctx context.Context, url string) 					(*dbmodels.Link, error)
}

type storageImpl struct {
	pool *pgxpool.Pool
	cfg  *config.Config
}

func NewstorageImpl(cfg *config.Config) (Storage, error) {
	
	databaseURL := fmt.Sprintf("%s:%s", cfg.DataBase.Host, cfg.DataBase.Port)
	
	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &storageImpl{
		pool: pool,
		cfg:  cfg,
	}, nil
}

func (s *storageImpl) GetLinkByURL(ctx context.Context, url string) (*dbmodels.Link, error) {

	var link dbmodels.Link

	query := `SELECT id, name FROM links WHERE name = $1`

	if err := s.pool.QueryRow(ctx, query, url).Scan(&link); err != nil {
		return nil, err
	}
	return &link, nil
}

func (s *storageImpl) GetURLS(ctx context.Context, userID uint) ([]dbmodels.Link, error) {
	var links []dbmodels.Link

	query := `
		SELECT l.id, l.url
		FROM links l
		JOIN link_users lu ON l.id = lu.link_id
		WHERE lu.user_id = $1
	`

	rows, err := s.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query links: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var link dbmodels.Link
		if err := rows.Scan(&link.ID, &link.Url); err != nil {
			return nil, fmt.Errorf("failed to scan link: %w", err)
		}
		links = append(links, link)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return links, nil
}

func (s *storageImpl) DeleteUser(ctx context.Context, userID uint) error {
	query := `DELETE FROM users WHERE user_id = $1`
	_, err := s.pool.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func (s *storageImpl) DeleteLink(ctx context.Context, linkID uint) error {
	query := `DELETE FROM links WHERE id = $1`
	_, err := s.pool.Exec(ctx, query, linkID)
	if err != nil {
		return fmt.Errorf("failed to delete link: %w", err)
	}
	return nil
}

func (s *storageImpl) DeleteLinkUser(ctx context.Context, linkID uint, userID uint) error {
	query := `DELETE FROM link_users WHERE link_id = $1 AND user_id = $2`
	_, err := s.pool.Exec(ctx, query, linkID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete link-user association: %w", err)
	}
	return nil
}

func (s *storageImpl) CreateLinkUser(ctx context.Context, linkID uint, userID uint) error {
	query := `INSERT INTO link_users (link_id, user_id) VALUES ($1, $2)`
	_, err := s.pool.Exec(ctx, query, linkID, userID)
	if err != nil {
		return fmt.Errorf("failed to create link-user association: %w", err)
	}
	return nil
}

func (s *storageImpl) CreateUser(ctx context.Context, userID uint, name string) error {
	query := `INSERT INTO users (user_id, name) VALUES ($1, $2)`
	_, err := s.pool.Exec(ctx, query, userID, name)
	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) && pgError.Code == ErrUniqueViolation {
			return nil
		}
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (s *storageImpl) CreateLink(ctx context.Context, link string) error {
	query := `INSERT INTO links (url) VALUES ($1)`
	_, err := s.pool.Exec(ctx, query, link)
	if err != nil {
		var pgError *pgconn.PgError 
		if errors.As(err, &pgError) && pgError.Code == ErrUniqueViolation {
			return nil 
		}
		return fmt.Errorf("failed to create link: %w", err)
	}
	return nil
}


func (s *storageImpl) GetLinkByID(ctx context.Context, id uint) (*dbmodels.Link, error) {
	var link dbmodels.Link
	query := `SELECT id, url FROM links WHERE id = $1`
	err := s.pool.QueryRow(ctx, query, id).Scan(&link.ID, &link.Url)
	if err != nil {
		return nil, fmt.Errorf("failed to get link by ID: %w", err)
	}
	return &link, nil
}