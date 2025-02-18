package models

import "gorm.io/gorm"

type LinkChat struct {
	gorm.Model
	LinkID uint `gorm:"primaryKey;column:link_id"`
	ChatID uint `gorm:"primaryKey;column:chat_id"`

	LinkRef Link `gorm:"foreignKey:LinkID;references:ID;constraint:OnDelete:CASCADE"`
	ChatRef Chat `gorm:"foreignKey:ChatID;references:ChatID;constraint:OnDelete:CASCADE"`
}
