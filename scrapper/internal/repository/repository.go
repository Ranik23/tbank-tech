package repository

import (
	"context"
	tx"tbank/scrapper/internal/repository/txmanager"
	dbmodels"tbank/scrapper/internal/models"
)





type Repository interface {
	tx.TxManager

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