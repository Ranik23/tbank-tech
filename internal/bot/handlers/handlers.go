package handlers

import "gopkg.in/telebot.v3"


type Handler interface {
	Handle(c telebot.Context) error
}