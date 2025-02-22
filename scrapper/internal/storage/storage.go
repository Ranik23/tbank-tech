package storage

import (
	"context"
	"errors"
	"tbank/scrapper/config"
	connection "tbank/scrapper/internal/db"
	dbmodels "tbank/scrapper/internal/db/models"

	"gorm.io/gorm"
)

type Storage interface {
	CreateChat(ctx context.Context, chatID int64)                         error
	CreateLinkChat(ctx context.Context, linkID int64, chatID int64)       error
	DeleteLinkChat(ctx context.Context, linkID int64, chatID int64)       error
	GetLinks(ctx context.Context, chatID int64)                           ([]dbmodels.Link, error)
	DeleteChat(ctx context.Context, chatID int64)                         error
	DeleteLink(ctx context.Context, linkID int64)                         error
	FindLinkID(ctx context.Context, link dbmodels.Link)                   (int64, error)
	CreateLink(ctx context.Context, link dbmodels.Link) 				  error
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


func (s *StorageImpl) CreateLink(ctx context.Context, link dbmodels.Link) error {
	tx := s.db.Begin()

	tx.Exec("SET TRANSACTION ISOLATION LEVEL REPEATABLE READ")

	if err := tx.Model(&link).Create(&link).Error; err != nil {
		defer tx.Rollback()
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil // не выводим ошибку
		}
		return err
	}
	return nil
}

func (s *StorageImpl) FindLinkID(ctx context.Context, link dbmodels.Link) (int64, error) {
	tx := s.db.Begin()

	tx.Exec("SET TRANSACTION ISOLATION LEVEL REPEATABLE READ")

	var l dbmodels.Link
	if err := tx.WithContext(ctx).Where("url = ?", link.Url).Find(&l).Error; err != nil {
		tx.Rollback()
		return -1, err
	}

	return int64(l.ID), tx.Commit().Error
}

func (s *StorageImpl) GetLinks(ctx context.Context, chatID int64) ([]dbmodels.Link, error) {
	tx := s.db.Begin()

	tx.Exec("SET TRANSACTION ISOLATION LEVEL REPEATABLE READ")

	if tx.Error != nil {
		return nil, tx.Error
	}

	var links []dbmodels.Link


	var linkChats []dbmodels.LinkChat
	err := tx.WithContext(ctx).Model(&linkChats).Where("chat_id = ?", chatID).Find(&linkChats).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	var linkIDs []int64
	for _, lc := range linkChats {
		linkIDs = append(linkIDs, int64(lc.LinkID))
	}

	if len(linkIDs) == 0 {
		tx.Rollback()
		return nil, nil
	}

	err = tx.Where("id IN (?)", linkIDs).Find(&links).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return links, tx.Commit().Error
}

func (s *StorageImpl) DeleteChat(ctx context.Context, chatID int64) error {
	tx := s.db.Begin()

	tx.Exec("SET TRANSACTION ISOLATION LEVEL REPEATABLE READ")

	if tx.Error != nil {
		return tx.Error
	}

	err := tx.WithContext(ctx).Where("id = ?", chatID).Delete(&dbmodels.Chat{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *StorageImpl) DeleteLink(ctx context.Context, linkID int64) error {
	tx := s.db.Begin()

	tx.Exec("SET TRANSACTION ISOLATION LEVEL REPEATABLE READ")

	if tx.Error != nil {
		return tx.Error
	}

	err := tx.WithContext(ctx).Where("id = ?", linkID).Delete(&dbmodels.Link{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *StorageImpl) DeleteLinkChat(ctx context.Context, linkID int64, chatID int64) error {
	tx := s.db.Begin()

	tx.Exec("SET TRANSACTION ISOLATION LEVEL REPEATABLE READ")

	if tx.Error != nil {
		return tx.Error
	}

	err := tx.WithContext(ctx).
		Where("link_id = ? AND chat_id = ?", linkID, chatID).
		Delete(&dbmodels.LinkChat{}).Error // ЧТОБЫ ДАННЫЕ КОТОРЫЕ МЫ ПРОЧИТАЛИ НЕ БЫЛИ ИЗМЕНЕНЫ ДРУГИМ, НАПРИМЕР УДАЛЕНЫ
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *StorageImpl) CreateLinkChat(ctx context.Context, linkID int64, chatID int64) error {
	tx := s.db.Begin()

	tx.Exec("SET TRANSACTION ISOLATION LEVEL READ COMMITTED")

	if tx.Error != nil {
		return tx.Error
	}

	LinkChat := dbmodels.LinkChat{
		LinkID: linkID,
		ChatID: chatID,
	}

	if err := tx.WithContext(ctx).Create(&LinkChat).Error; err != nil {
		defer tx.Rollback()
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil
		}
		return err
	}

	return tx.Commit().Error
}

func (s *StorageImpl) CreateChat(ctx context.Context, chatID int64) error {
	tx := s.db.Begin()

	tx.Exec("SET TRANSACTION ISOLATION LEVEL READ COMMITTED")

	if tx.Error != nil {
		return tx.Error
	}

	chat := dbmodels.Chat{ChatID: chatID}
	if err := tx.WithContext(ctx).Create(&chat).Error; err != nil {
		defer tx.Rollback()
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil
		}
		return err
	}

	return tx.Commit().Error
}
