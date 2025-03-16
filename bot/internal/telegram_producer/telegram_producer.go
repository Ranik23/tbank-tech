package telegramproducer

import (
	"log/slog"
	telegrambot "tbank/bot/internal/telegram_bot"
	"tbank/bot/internal/telegram_producer/utils"

	"github.com/IBM/sarama"
	"github.com/google/go-github/v69/github"
	"gopkg.in/telebot.v3"
)


type TelegramProducer struct {
	bot 		telegrambot.TelegramBot
	commitCh 	chan sarama.ConsumerMessage
	stopCh		chan struct{}
	logger		*slog.Logger
	workerDone 	chan struct{}
}


func NewTelegramProducer(telegramBot telegrambot.TelegramBot, logger *slog.Logger, commitCh chan sarama.ConsumerMessage) *TelegramProducer {
	return &TelegramProducer{
		bot: 		telegramBot,
		commitCh: 	commitCh,
		stopCh: 	make(chan struct{}),
		logger: 	logger,
		workerDone: make(chan struct{}),
	}
}

func (tp *TelegramProducer) Run() {
	const op = "TelegramProducer.Run"
	tp.logger.Info(op, slog.String("msg","Telegram Producer"))
	go func() {
		defer close(tp.workerDone)

		for {
			select {
			case message := <-tp.commitCh:
				tp.logger.Info("op", slog.String("msg", "got the message"))
				commit, err := utils.ConvertFromBytesToCustomCommit(message.Value)
				if err != nil {
					tp.logger.Error(op, slog.String("error", err.Error()))
					tp.Stop()
				}

				userID := commit.UserID
				comm := commit.Commit

				if err := tp.sendMessageToUser(int64(userID), comm); err != nil {
					tp.logger.Error(op, slog.String("error", err.Error()))
				}
			case <-tp.stopCh:
				tp.logger.Error(op, slog.String("msg", "stopping the producer"))
				return
			}
		}
	}()
}


func (tp *TelegramProducer) sendMessageToUser(userID int64, commit *github.RepositoryCommit) error {
	const op = "TelegramProducer.sendMessageToUser"
	tp.logger.Info(op, slog.String("msg", "sending the commit"))
	_, err := tp.bot.Send(&telebot.User{ID: userID}, commit)
	if err != nil {
		tp.logger.Error(op, "error", err.Error())
		return err
	}
	return nil
}


func (tp *TelegramProducer) Stop() {
	close(tp.stopCh)
	<-tp.workerDone
}