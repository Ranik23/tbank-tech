package models

import (
	"github.com/google/go-github/v69/github"
)


type CustomCommit struct {
	Commit *github.RepositoryCommit	`json:"commit"`
	UserID uint						`json:"user_id"`
}

type User struct {
	State   int
	Link    string
	Name	string
	Token	string
}