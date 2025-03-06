package models

import (
	"gorm.io/gorm"
)



// Оставляем тут UserID потому что мы берем user_id из бота чата
type User struct {
	gorm.Model
	UserID	uint		`gorm:"primaryKey;column:user_id"` 
	Name	string		`gorm:"column:name"`
}