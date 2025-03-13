package telegramproducer

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"gopkg.in/telebot.v3"
)



type TelegramProducer struct {
	bot *telebot.Bot
	messageCh chan kafka.Message
}


func NewTelegramProducer(telegramBot *telebot.Bot, messagesCh chan kafka.Message) *TelegramProducer {
	return &TelegramProducer{
		bot: telegramBot,
		messageCh: messagesCh,
	}
}

func (tp *TelegramProducer) Run() {
	return
}



func (tp *TelegramProducer) sendMessageToUser(msg *kafka.Message) error {
	return nil
}

func (tp *TelegramProducer) Stop() {
	return
}