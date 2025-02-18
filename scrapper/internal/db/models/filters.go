package models

import "gorm.io/gorm"


type Filter struct {
	gorm.Model
	Name string		`gorm:"primaryKey;not null"`
}