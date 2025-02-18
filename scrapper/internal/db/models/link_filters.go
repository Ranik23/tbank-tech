package models

import "gorm.io/gorm"


type LinkFilters struct {
	gorm.Model
	LinkID 		uint		`gorm:"primaryKey"`
	FilterID	uint		`gorm:"primaryKey"`

	LinkRef		Link 		`gorm:"foreignKey:LinkID;references:ID;constraint:OnDelete:CASCADE"`
	FilterRef	Filter		`gorm:"foreignKey:FilterID;references:ID;constraint:OnDelete:CASCADE"`
}