package storage

import (
	"context"
	"tbank/scrapper/config"
	connection "tbank/scrapper/internal/db"
	dbmodels "tbank/scrapper/internal/db/models"

	"gorm.io/gorm"
)

type Storage interface {
	CreateChat(ctx context.Context, chatID uint) 						error
	CreateFilter(ctx context.Context, name string) 						error
	CreateTag(ctx context.Context, name string) 						error
	CreateLinkChat(ctx context.Context, linkID uint, chatID uint) 		error
	CreateLinkFilter(ctx context.Context, linkID uint, filterID uint) 	error
	CreateLinkTag(ctx context.Context, linkID uint, tagID uint) 		error

	DeleteChat(ctx context.Context, chatID uint) 						error
	DeleteLink(ctx context.Context, linkID uint) 						error 
	DeleteLinkChat(ctx context.Context, linkID uint, chatID uint) 		error

	GetURLS(ctx context.Context, chatID uint) 							([]dbmodels.Link, error)
}

type StorageImpl struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewStorageImpl(cfg *config.Config) (*StorageImpl, error) {
	db, err := connection.ConnectToDB(cfg)
	if err != nil {
		return nil, err
	}
	return &StorageImpl{
		db:  db,
		cfg: cfg,
	}, nil
}

func (s *StorageImpl) GetURLs(ctx context.Context, chatID uint) ([]dbmodels.Link, error) {
	var links []dbmodels.Link

	err := s.db.WithContext(ctx).
		Joins("JOIN link_chats ON link_chats.link_id = links.id").
		Where("link_chats.chat_id = ?", chatID).
		Find(&links).Error

	if err != nil {
		return nil, err
	}

	return links, nil
}

func (s *StorageImpl) DeleteChat(ctx context.Context, chatID uint) error {
	return s.db.WithContext(ctx).Where("id = ?", chatID).Delete(&dbmodels.Chat{}).Error
}

func (s *StorageImpl) DeleteLink(ctx context.Context, linkID uint) error {
	return s.db.WithContext(ctx).Where("id = ?", linkID).Delete(&dbmodels.Link{}).Error
}

func (s *StorageImpl) DeleteLinkChat(ctx context.Context, linkID uint, chatID uint) error {
	return s.db.WithContext(ctx).
		Where("link_id = ? AND chat_id = ?", linkID, chatID).
		Delete(&dbmodels.LinkChat{}).Error
}

func (s *StorageImpl) CreateLinkChat(ctx context.Context, linkID uint, chatID uint) error {
	LinkChat := dbmodels.LinkChat{
		LinkID: linkID,
		ChatID: chatID,
	}
	return s.db.WithContext(ctx).Create(&LinkChat).Error
}

func (s *StorageImpl) CreateChat(ctx context.Context, chatID uint) error {
	chat := dbmodels.Chat{ChatID: chatID}
	return s.db.WithContext(ctx).Create(&chat).Error
}

func (s *StorageImpl) CreateFilter(ctx context.Context, name string) error {
	filter := dbmodels.Filter{Name: name}
	return s.db.WithContext(ctx).Create(&filter).Error
}

func (s *StorageImpl) CreateTag(ctx context.Context, name string) error {
	tag := dbmodels.Tag{Name: name}
	return s.db.WithContext(ctx).Create(&tag).Error
}

func (s *StorageImpl) CreateLinkFilter(ctx context.Context, linkID uint, filterID uint) error {
	linkFilter := dbmodels.LinkFilters{LinkID: linkID, FilterID: filterID}
	return s.db.WithContext(ctx).Create(&linkFilter).Error
}

func (s *StorageImpl) CreateLinkTag(ctx context.Context, linkID uint, tagID uint) error {
	linkTag := dbmodels.LinkTags{LinkID: linkID, TagID: tagID}
	return s.db.WithContext(ctx).Create(&linkTag).Error
}
