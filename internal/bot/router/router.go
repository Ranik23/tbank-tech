package router

import (
	"fmt"
	"tbank/internal/bot/handlers"
	"tbank/internal/usecase"

	"gopkg.in/telebot.v3"
)

var (
	ErrPathAlreadyExists = fmt.Errorf("path already exists")
)


type Router struct {
	bot 		*telebot.Bot
	handlers 	map[string]handlers.Handler
}

func NewRouter(bot *telebot.Bot, usecase usecase.UseCase) *Router {
	return &Router{
		bot: bot,
		handlers: make(map[string]handlers.Handler),
	}
}

func (r *Router) RegisterHandlers() {
	for command, handler := range r.handlers {
		r.bot.Handle(command, handler.Handle)
	}
}

func (r *Router) AddHandler(path string, handler handlers.Handler) error {
	_, ok := r.handlers[path]
	if ok {
		return ErrPathAlreadyExists
	}
	r.handlers[path] = handler
	return nil
}