//go:build unit

package kafkaconsumer

import (
	"fmt"
	"log/slog"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/IBM/sarama/mocks"
	"github.com/stretchr/testify/require"
)

func TestKafkaConsumer_Success(t *testing.T) {
	exampleTopic := "test"
	examplePartitions := []int32{1, 2, 3}

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	mockConsumer := mocks.NewConsumer(t, config)

	mockConsumer.SetTopicMetadata(map[string][]int32{
		exampleTopic: examplePartitions,
	})

	logger := slog.Default()

	mockPartitionConsumer1 := mockConsumer.ExpectConsumePartition(exampleTopic, examplePartitions[0], sarama.OffsetNewest)
	mockPartitionConsumer2 := mockConsumer.ExpectConsumePartition(exampleTopic, examplePartitions[1], sarama.OffsetNewest)
	mockPartitionConsumer3 := mockConsumer.ExpectConsumePartition(exampleTopic, examplePartitions[2], sarama.OffsetNewest)

	// Отправляем три тестовых сообщения
	messages := []string{"message1", "message2", "message3"}

	mockPartitionConsumer1.YieldMessage(&sarama.ConsumerMessage{Topic: exampleTopic, Value: []byte(messages[0])})
	mockPartitionConsumer2.YieldMessage(&sarama.ConsumerMessage{Topic: exampleTopic, Value: []byte(messages[1])})
	mockPartitionConsumer3.YieldMessage(&sarama.ConsumerMessage{Topic: exampleTopic, Value: []byte(messages[2])})

	commitCh := make(chan sarama.ConsumerMessage)
	kafkaConsumer := &KafkaConsumer{
		consumer: mockConsumer,
		commitCh: commitCh,
		stopCh: make(chan struct{}),
		logger: logger,
		topic: exampleTopic,
	}


	kafkaConsumer.Run()

	defer kafkaConsumer.Stop()

	errorCh := make(chan error, 1)
	receivedMessages := make(map[string]bool)

	go func() {
		for commit := range commitCh {
			receivedMessages[string(commit.Value)] = true
			if commit.Topic != exampleTopic {
				errorCh <- fmt.Errorf("Topic mismatch: got %s, expected %s", commit.Topic, exampleTopic)
				return
			}
		}
	}()

	select {
	case err := <-errorCh:
		t.Fatalf("Failed: %v", err)
	case <-time.After(2 * time.Second):
		for _, msg := range messages {
			require.True(t, receivedMessages[msg], "Message not received: %s", msg)
		}
	}
}

