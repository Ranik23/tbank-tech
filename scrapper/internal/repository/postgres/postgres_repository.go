package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	dbmodels "github.com/Ranik23/tbank-tech/scrapper/internal/models"
	"github.com/Ranik23/tbank-tech/scrapper/internal/repository"

	"github.com/jackc/pgx/v5"
)

var (
	ErrNoUserFound           = errors.New("no user found")
	ErrNoLinkUserFound       = errors.New("no link-user found")
	ErrFailedToGetUser       = errors.New("failed to get user")
	ErrNoUsersFound          = errors.New("no users found")
	ErrNoLinkFound           = errors.New("no link found")
	ErrFailedToQuery         = errors.New("failed to query")
	ErrFailedToScan          = errors.New("failed to scan")
	ErrRowsError             = errors.New("rows error")
	ErrFailedToUpdate        = errors.New("failed to update")
	ErrFailedToDelete        = errors.New("failed to delete")
	ErrFailedToCreate        = errors.New("failed to create")
	ErrNoRowsAffected        = errors.New("no rows affected")
	ErrLinkAlreadyExists     = errors.New("link already exists")
	ErrUserAlreadyExists     = errors.New("user already exists")
	ErrLinkUserAlreadyExists = errors.New("linkuser already exists")
	ErrFailedToGetLinkUser   = errors.New("failed to get linkuser")
)

type postgresRepository struct {
	logger    *slog.Logger
	txManager repository.TxManager
}

func NewPostgresRepository(txManager repository.TxManager, logger *slog.Logger) repository.Repository {
	return &postgresRepository{txManager: txManager, logger: logger}
}

func (s *postgresRepository) GetLinkUser(ctx context.Context, userID uint, linkID uint) (linkuser *dbmodels.LinkUser, err error) {
	s.logger.Info("GetLinkUser called", slog.Int64("userID", int64(userID)), slog.Int64("linkID", int64(linkID)))

	executor := s.txManager.GetExecutor(ctx)

	query := `SELECT link_id, user_id FROM link_users WHERE link_id = $1 AND user_id = $2`

	linkuser = &dbmodels.LinkUser{}

	err = executor.QueryRow(ctx, query, userID, linkID).Scan(&linkuser.LinkID, &linkuser.UserID)
	if err != nil {
		s.logger.Error("GetLinkUser failed", slog.String("erorr", err.Error()))
		return nil, fmt.Errorf("%w: %v", ErrFailedToGetLinkUser, err)
	}

	s.logger.Info("GetLinkUser success", slog.Int64("userID", int64(userID)), slog.Int64("linkID", int64(linkID)))
	return linkuser, nil
}

// error if operation failes or not found
func (s *postgresRepository) GetUserByName(ctx context.Context, name string) (user *dbmodels.User, err error) {
	s.logger.Info("GetUserByName called", slog.String("name", name))

	executor := s.txManager.GetExecutor(ctx)

	query := `SELECT user_id, name FROM users WHERE name = $1`

	user = &dbmodels.User{}
	err = executor.QueryRow(ctx, query, name).Scan(&user.UserID, &user.Name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.logger.Warn("No user found with", slog.String("name", name))
			return nil, ErrNoUserFound
		}
		s.logger.Error("GetUserByName failed", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%w: %v", ErrFailedToGetUser, err)
	}

	s.logger.Info("GetUserByName success", slog.String("name", name))
	return user, nil
}


// nil, NoUserFound - если такого пользователя нет, (nil, err) - ошибка поиска, (user, nil) - нашли
func (s *postgresRepository) GetUserByID(ctx context.Context, userID uint) (user *dbmodels.User, err error) {
	s.logger.Info("GetUserByID called", slog.Int64("userID", int64(userID)))

	executor := s.txManager.GetExecutor(ctx)

	query := `SELECT user_id, name, token FROM users WHERE user_id = $1`

	user = &dbmodels.User{}
	err = executor.QueryRow(ctx, query, userID).Scan(&user.UserID, &user.Name, &user.Token)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.logger.Warn("No user found with", slog.Int64("userID", int64(userID)))
			return nil, ErrNoUserFound
		}
		s.logger.Error("GetUserByID failed", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%w: %v", ErrFailedToGetUser, err)
	}
	s.logger.Info("GetUserByID success", slog.Int64("userID", int64(userID)))

	return user, nil
}

func (s *postgresRepository) GetUsers(ctx context.Context) (users []dbmodels.User, err error) {
	s.logger.Info("GetUser called")

	executor := s.txManager.GetExecutor(ctx)

	query := `SELECT user_id, name FROM users`

	rows, err := executor.Query(ctx, query)
	if err != nil {
		s.logger.Error("GetUsers failed to query", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%w: %v", ErrFailedToQuery, err)
	}
	defer rows.Close()

	for rows.Next() {
		var user dbmodels.User
		if err := rows.Scan(&user.UserID, &user.Name); err != nil {
			s.logger.Error("GetUsers failed to scan row", slog.String("error", err.Error()))
			return nil, fmt.Errorf("%w: %v", ErrFailedToScan, err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		s.logger.Error("GetUsers rows error", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%w: %v", ErrRowsError, err)
	}

	return users, nil
}

func (s *postgresRepository) GetLinkByURL(ctx context.Context, url string) (link *dbmodels.Link, err error) {
	s.logger.Info("GetLink called", slog.String("url", url))

	executor := s.txManager.GetExecutor(ctx)
	query := `SELECT id, url FROM links WHERE url = $1`

	link = &dbmodels.Link{}
	err = executor.QueryRow(ctx, query, url).Scan(&link.ID, &link.Url)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.logger.Warn("No link found with", slog.String("link", url))
			return nil, ErrNoLinkFound
		}
		s.logger.Error("GetLink failed", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%w: %v", ErrFailedToQuery, err)
	}
	s.logger.Info("GetLinkByURL success", slog.String("url", url), slog.Uint64("linkID", uint64(link.ID)))
	return link, nil
}

func (s *postgresRepository) GetLinks(ctx context.Context, userID uint) (links []dbmodels.Link, err error) {
	s.logger.Info("GetLinks called", slog.Uint64("userID", uint64(userID)))

	executor := s.txManager.GetExecutor(ctx)

	query := `
		SELECT l.id, l.url
		FROM links l
		JOIN link_users lu ON l.id = lu.link_id
		WHERE lu.user_id = $1
	`

	rows, err := executor.Query(ctx, query, userID)
	if err != nil {
		s.logger.Error("GetLinks failed", slog.String("error", err.Error()))
		return nil, ErrFailedToQuery
	}
	defer rows.Close()

	for rows.Next() {
		var link dbmodels.Link
		if err := rows.Scan(&link.ID, &link.Url); err != nil {
			s.logger.Error("GetLinks scan failed", slog.String("error", err.Error()))
			return nil, fmt.Errorf("%w: %v", ErrFailedToScan, err)
		}
		links = append(links, link)
	}

	if err := rows.Err(); err != nil {
		s.logger.Error("GetLinks rows error", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%w: %v", ErrRowsError, err)
	}

	s.logger.Info("GetLinks success", slog.Int("count", len(links)))
	return links, nil
}

func (s *postgresRepository) DeleteUser(ctx context.Context, userID uint) error {
	s.logger.Info("DeleteUser called", slog.Uint64("userID", uint64(userID)))

	executor := s.txManager.GetExecutor(ctx)
	query := `DELETE FROM users WHERE user_id = $1`

	result, err := executor.Exec(ctx, query, userID)
	if err != nil {
		s.logger.Error("DeleteUser failed", slog.String("error", err.Error()))
		return fmt.Errorf("%w: %v", ErrFailedToDelete, err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		s.logger.Warn("DeleteUser: user not found", slog.Uint64("userID", uint64(userID)))
		return fmt.Errorf("%w: %v", ErrNoUserFound, err)
	}

	s.logger.Info("DeleteUser success", slog.Int64("rowsAffected", rowsAffected))
	return nil
}

func (s *postgresRepository) DeleteLink(ctx context.Context, linkID uint) error {
	s.logger.Info("DeleteLink called", slog.Uint64("linkID", uint64(linkID)))

	executor := s.txManager.GetExecutor(ctx)
	query := `DELETE FROM links WHERE id = $1`

	result, err := executor.Exec(ctx, query, linkID)
	if err != nil {
		s.logger.Error("DeleteLink failed", slog.String("error", err.Error()))
		return fmt.Errorf("%w: %v", ErrFailedToDelete, err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		s.logger.Warn("DeleteLink: link not found", slog.Uint64("linkID", uint64(linkID)))
		return fmt.Errorf("%w: %v", ErrNoLinkFound, err)
	}

	s.logger.Info("DeleteLink success", slog.Int64("rowsAffected", rowsAffected))
	return nil
}

func (s *postgresRepository) DeleteLinkUser(ctx context.Context, linkID uint, userID uint) error {
	s.logger.Info("DeleteLinkUser called",
		slog.Uint64("linkID", uint64(linkID)),
		slog.Uint64("userID", uint64(userID)),
	)

	executor := s.txManager.GetExecutor(ctx)
	query := `DELETE FROM link_users WHERE link_id = $1 AND user_id = $2`

	result, err := executor.Exec(ctx, query, linkID, userID)
	if err != nil {
		s.logger.Error("DeleteLinkUser failed",
			slog.String("error", err.Error()),
			slog.Uint64("linkID", uint64(linkID)),
			slog.Uint64("userID", uint64(userID)),
		)
		return fmt.Errorf("%w: %v", ErrFailedToDelete, err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		s.logger.Warn("DeleteLinkUser: link-user association not found",
			slog.Uint64("linkID", uint64(linkID)),
			slog.Uint64("userID", uint64(userID)),
		)
		return fmt.Errorf("%w: %v", ErrNoLinkUserFound, err)
	}

	s.logger.Info("DeleteLinkUser success", slog.Int64("rowsAffected", rowsAffected))
	return nil
}

func (s *postgresRepository) CreateLinkUser(ctx context.Context, linkID uint, userID uint) error {
	s.logger.Info("CreateLinkUser called",
		slog.Uint64("linkID", uint64(linkID)),
		slog.Uint64("userID", uint64(userID)),
	)

	executor := s.txManager.GetExecutor(ctx)
	query := `INSERT INTO link_users (link_id, user_id) VALUES ($1, $2)`

	_, err := executor.Exec(ctx, query, linkID, userID)
	if err != nil {
		s.logger.Error("CreateLinkUser failed",
			slog.String("error", err.Error()),
			slog.Uint64("linkID", uint64(linkID)),
			slog.Uint64("userID", uint64(userID)),
		)
		return fmt.Errorf("%w: %v", ErrFailedToCreate, err)
	}

	s.logger.Info("CreateLinkUser success",
		slog.Uint64("linkID", uint64(linkID)),
		slog.Uint64("userID", uint64(userID)),
	)
	return nil
}

func (s *postgresRepository) CreateUser(ctx context.Context, userID uint, name string, token string) error {
	s.logger.Info("CreateUser called", slog.Uint64("userID", uint64(userID)), slog.String("name", name))

	executor := s.txManager.GetExecutor(ctx)
	query := `INSERT INTO users (user_id, name, token) VALUES ($1, $2, $3)`

	_, err := executor.Exec(ctx, query, userID, name, token)
	if err != nil {
		s.logger.Error("CreateUser failed", slog.String("error", err.Error()))
		return fmt.Errorf("%w: %v", ErrFailedToCreate, err)
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
		s.logger.Error("CreateLink failed", slog.String("error", err.Error()))
		return fmt.Errorf("%w: %v", ErrFailedToCreate, err)
	}

	s.logger.Info("CreateLink success")
	return nil
}
