package models

import "gorm.io/gorm"


type LinkTags struct {
	gorm.Model
	LinkID	uint	`gorm:"primaryKey"`
	TagID	uint	`gorm:"primaryKey"`
	LinkRef Link 	`gorm:"foreignKey:LinkID;references:ID;constraint:OnDelete:CASCADE"`
	TagRef 	Tag 	`gorm:"foreignKey:TagID;references:ID;constraint:OnDelete:CASCADE"`
}