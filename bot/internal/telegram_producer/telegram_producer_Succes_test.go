package telegramproducer

import (
	"encoding/json"
	"log/slog"
	"sync"
	"testing"
	"time"

	"github.com/Ranik23/tbank-tech/bot/internal/models"
	mockbot "github.com/Ranik23/tbank-tech/bot/internal/telegram_bot/mock"

	"github.com/IBM/sarama"
	"github.com/golang/mock/gomock"
	"github.com/google/go-github/v69/github"
	"github.com/stretchr/testify/require"
	"gopkg.in/telebot.v3"
)

func TestTelegramProducer_Success(t *testing.T) {

	exampleCommit := models.CustomCommit{
		UserID: 1,
		Commit: &github.RepositoryCommit{
			SHA: github.Ptr("test_sha"),
			Commit: &github.Commit{
				Message: github.Ptr("test_message"),
			},
		},
	}

	messageCh := make(chan sarama.ConsumerMessage)

	ctrl := gomock.NewController(t)
	mockBot := mockbot.NewMockTelegramBot(ctrl)

	mockBot.EXPECT().Send(gomock.Eq(&telebot.User{ID: 1}), gomock.Any(), gomock.Any()).Times(1).Return(nil, nil)

	producerTelegram := NewTelegramProducer(mockBot, slog.Default(), messageCh)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		producerTelegram.Run()
	}()

	jsonExampleCommit, err := json.Marshal(exampleCommit)
	require.NoError(t, err)

	select {
	case messageCh <- sarama.ConsumerMessage{Value: jsonExampleCommit}:
	case <-time.After(5 * time.Second):
		t.Fatalf("Timeout Expired")
	}

	time.Sleep(5 * time.Second)
	producerTelegram.Stop()
	wg.Wait()
}
