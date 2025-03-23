//go:build unit

package kafkaconsumer

import (
	"log/slog"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/IBM/sarama/mocks"
	"github.com/stretchr/testify/require"
)




func TestKafkaConsumer_DuplicateMessages(t *testing.T) {
	exampleTopic := "test"
	examplePartitions := []int32{0}

	logger := slog.Default()

	config := sarama.NewConfig()
	mockConsumer := mocks.NewConsumer(t, config)
	mockConsumer.SetTopicMetadata(map[string][]int32{exampleTopic: examplePartitions})

	mockPartitionConsumer := mockConsumer.ExpectConsumePartition(exampleTopic, examplePartitions[0], sarama.OffsetNewest)

	message := &sarama.ConsumerMessage{Topic: exampleTopic, Value: []byte("duplicate")}
	mockPartitionConsumer.YieldMessage(message)
	mockPartitionConsumer.YieldMessage(message)

	commitCh := make(chan sarama.ConsumerMessage, 2)
	kafkaConsumer:= &KafkaConsumer{
		consumer: mockConsumer,
		commitCh: commitCh,
		stopCh: make(chan struct{}),
		logger: logger,
		topic: exampleTopic,
	}

	kafkaConsumer.Run()
	defer kafkaConsumer.Stop()

	received := make(map[string]int)

	timeout := time.After(2 * time.Second)

	for i := 0; i < 2; i++ {
		select {
		case msg := <-commitCh:
			received[string(msg.Value)]++
		case <-timeout:
			t.Fatal("Timeout waiting for messages")
		}
	}

	require.Equal(t, 2, received["duplicate"], "Duplicate message not received twice")
}
