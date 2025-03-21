//go:build unit

package telegramproducer

import (
	"encoding/json"
	"log/slog"
	"sync"
	"testing"
	"time"

	"github.com/Ranik23/tbank-tech/bot/internal/models"
	botmock "github.com/Ranik23/tbank-tech/bot/internal/telegram_producer/telegram_bot/mock"

	"github.com/IBM/sarama"
	"github.com/golang/mock/gomock"
	"github.com/google/go-github/v69/github"
	"github.com/stretchr/testify/require"
	"gopkg.in/telebot.v3"
)

func TestTelegramProducer_MultipleMessages(t *testing.T) {
	exampleCommit1 := models.CustomCommit{
		UserID: 1,
		Commit: &github.RepositoryCommit{
			SHA: github.Ptr("sha_1"),
			Commit: &github.Commit{
				Message: github.Ptr("message_1"),
			},
		},
	}

	exampleCommit2 := models.CustomCommit{
		UserID: 2,
		Commit: &github.RepositoryCommit{
			SHA: github.Ptr("sha_2"),
			Commit: &github.Commit{
				Message: github.Ptr("message_2"),
			},
		},
	}

	messageCh := make(chan sarama.ConsumerMessage)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBot := botmock.NewMockTelegramBot(ctrl)

	mockBot.EXPECT().Send(gomock.Eq(&telebot.User{ID: 1}), gomock.Any(), gomock.Any()).Times(1).Return(nil, nil)
	mockBot.EXPECT().Send(gomock.Eq(&telebot.User{ID: 2}), gomock.Any(), gomock.Any()).Times(1).Return(nil, nil)

	producerTelegram := NewTelegramProducer(mockBot, slog.Default(), messageCh)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		producerTelegram.Run()
	}()

	jsonCommit1, err := json.Marshal(exampleCommit1)
	require.NoError(t, err)
	jsonCommit2, err := json.Marshal(exampleCommit2)
	require.NoError(t, err)

	messageCh <- sarama.ConsumerMessage{Value: jsonCommit1}
	messageCh <- sarama.ConsumerMessage{Value: jsonCommit2}

	time.Sleep(1 * time.Second)
	producerTelegram.Stop()
	wg.Wait()
}
