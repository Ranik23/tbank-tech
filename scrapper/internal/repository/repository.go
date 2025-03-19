package repository

import (
	"context"

	dbmodels "github.com/Ranik23/tbank-tech/scrapper/internal/models"
)

type Repository interface {
	CreateLink(ctx context.Context, link string) error
	CreateUser(ctx context.Context, userID uint, name string, token string) error
	CreateLinkUser(ctx context.Context, linkID uint, userID uint) error

	DeleteUser(ctx context.Context, userID uint) error
	DeleteLink(ctx context.Context, linkID uint) error
	DeleteLinkUser(ctx context.Context, linkID uint, userID uint) error

	GetLinks(ctx context.Context, userID uint) (links []dbmodels.Link, err error)
	GetLinkByURL(ctx context.Context, url string) (link *dbmodels.Link, err error)
	GetUserByName(ctx context.Context, name string) (user *dbmodels.User, err error)
	GetUserByID(ctx context.Context, userID uint) (user *dbmodels.User, err error)
	GetUsers(ctx context.Context) (users []dbmodels.User, err error)
	GetLinkUser(ctx context.Context, userID uint, linkID uint) (linkuser *dbmodels.LinkUser, err error)
}
