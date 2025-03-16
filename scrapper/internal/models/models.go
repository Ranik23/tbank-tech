package models

import "gorm.io/gorm"

type LinkUser struct {
	gorm.Model
	LinkID uint `gorm:"primaryKey;column:link_id"`
	UserID uint `gorm:"primaryKey;column:user_id"`

	LinkRef Link `gorm:"foreignKey:LinkID;references:ID;constraint:OnDelete:CASCADE"`
	ChatRef User `gorm:"foreignKey:UserID;references:UserID;constraint:OnDelete:CASCADE"`
}

type Link struct {
	gorm.Model
	Url string		`gorm:"column:url;not null"`
}

type User struct {
	gorm.Model
	UserID	uint		`gorm:"primaryKey;column:user_id"` 
	Name	string		`gorm:"column:name"`
}