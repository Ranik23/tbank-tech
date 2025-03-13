package telegramproducer

import (
	"encoding/json"
	"log/slog"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/go-github/v69/github"
	"gopkg.in/telebot.v3"
)

type CustomCommit struct {
	Commit *github.RepositoryCommit	`json:"commit"`
	UserID uint						`json:"user_id"`
}


type TelegramProducer struct {
	bot 		*telebot.Bot
	messageCh 	chan kafka.Message
	stopCh		chan struct{}
	logger		*slog.Logger
}


func NewTelegramProducer(telegramBot *telebot.Bot, logger *slog.Logger, messagesCh chan kafka.Message) *TelegramProducer {
	return &TelegramProducer{
		bot: telegramBot,
		messageCh: messagesCh,
		stopCh: make(chan struct{}),
		logger: logger,
	}
}

func (tp *TelegramProducer) Run() {
	const op = "TelegramProducer.Run"
	tp.logger.Info(op, slog.String("msg","Telegram Producer"))
	go func() {
		for {
			select {
			case message := <-tp.messageCh:
				commit, err := tp.convert(message.Value)
				if err != nil {
					tp.logger.Error(op, slog.String("error", err.Error()))
					continue
				}

				userID := commit.UserID
				comm := commit.Commit

				if err := tp.sendMessageToUser(comm, userID); err != nil {
					tp.logger.Error(op, slog.String("error", err.Error()))
				}
			case <-tp.stopCh:
				tp.logger.Error(op, slog.String("stopped", "true"))
				return
			}
		}
	}()
}


func (tp *TelegramProducer) convert(message []byte) (*CustomCommit, error) {
	const op = "TelegramProducer.convert"
	var msg CustomCommit
	if err := json.Unmarshal(message, &msg); err != nil {
		tp.logger.Error(op, slog.String("error", err.Error()))
		return nil, err
	}
	return &msg, nil
}


func (tp *TelegramProducer) sendMessageToUser(commit *github.RepositoryCommit, userID uint) error {
	const op = "TelegramProducer.sendMessageToUser"
	_, err := tp.bot.Send(&telebot.User{ID: int64(userID)}, commit)
	if err != nil {
		tp.logger.Error(op, "error", err.Error())
		return err
	}
	return nil
}

func (tp *TelegramProducer) Stop() {
	close(tp.stopCh)
}