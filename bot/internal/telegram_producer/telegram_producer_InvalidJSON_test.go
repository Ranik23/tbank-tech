package telegramproducer

import (
	"log/slog"
	"sync"
	"testing"
	"time"

	botmock "github.com/Ranik23/tbank-tech/bot/internal/telegram_bot/mock"

	"github.com/IBM/sarama"
	"github.com/golang/mock/gomock"
)

func TestTelegramProducer_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	messageCh := make(chan sarama.ConsumerMessage, 1)
	mockBot := botmock.NewMockTelegramBot(ctrl)

	mockBot.EXPECT().Send(gomock.Any(), gomock.Any()).Times(0)

	producerTelegram := NewTelegramProducer(mockBot, slog.Default(), messageCh)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		producerTelegram.Run()
	}()

	select {
	case messageCh <- sarama.ConsumerMessage{Value: sarama.ByteEncoder("invalid_json")}:
	case <-time.After(150 * time.Millisecond):
		t.Fatalf("Timeout Expired")
	}

	producerTelegram.Stop()
	wg.Wait()
}
