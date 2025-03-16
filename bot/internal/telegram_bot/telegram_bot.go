package telegrambot

import "gopkg.in/telebot.v3"


type TelegramBot interface {
	Send(to telebot.Recipient, what interface{}, opts ...interface{}) (*telebot.Message, error) 
}
