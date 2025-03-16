package models

type LinkUser struct {
	LinkID uint 
	UserID uint 
	LinkRef Link 
	ChatRef User
}

type Link struct {
	ID	uint
	Url string
}

type User struct {
	UserID	uint
	Name	string
}