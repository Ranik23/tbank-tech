package db

import (
	"fmt"
	"log/slog"
	"tbank/scrapper/config"
	"tbank/scrapper/internal/db/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


func ConnectToDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DataBase.Username, cfg.DataBase.Password, cfg.DataBase.Host, cfg.DataBase.Port, cfg.DataBase.DBName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.Chat{},
			&models.Filter{},
			&models.Link{},
			&models.LinkFilters{},
			&models.LinkTags{},
			&models.Tag{},
			&models.LinkChat{})

	if err != nil {
		return nil, err
	}

	slog.Info("Cpnnected to DB!")
	return db, nil
}