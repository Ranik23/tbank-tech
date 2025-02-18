package models

import (
	"gorm.io/gorm"
)



type Chat struct {
	gorm.Model
	ChatID	uint		`gorm:"primaryKey;column:chat_id"`
}