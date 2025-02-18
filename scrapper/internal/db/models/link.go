package models

import "gorm.io/gorm"


type Link struct {
	gorm.Model
	Url string		`gorm:"column:url;not null"`
}