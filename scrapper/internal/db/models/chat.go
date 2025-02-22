package models

import (
	"gorm.io/gorm"
)



type Chat struct {
	gorm.Model
	ChatID	int64		`gorm:"primaryKey;column:chat_id"`
}