//go:build unit

package kafkaconsumer

import (
	"log/slog"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/IBM/sarama/mocks"
)




func TestKafkaConsumer_NoMessages(t *testing.T) {
	exampleTopic := "test"
	examplePartitions := []int32{0}

	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Return.Errors = true

	mockConsumer := mocks.NewConsumer(t, saramaConfig)

	mockConsumer.SetTopicMetadata(map[string][]int32{exampleTopic: examplePartitions})
	mockConsumer.ExpectConsumePartition(exampleTopic, examplePartitions[0], sarama.OffsetNewest)

	commitCh := make(chan sarama.ConsumerMessage)
	kafkaConsumer := KafkaConsumer{
		consumer: mockConsumer,
		commitCh: commitCh,
		stopCh: make(chan struct{}),
		logger: slog.Default(),
		topic: exampleTopic,
	}


	kafkaConsumer.Run()
	defer kafkaConsumer.Stop()

	select {
	case <-commitCh:
		t.Fatal("Unexpected message received")
	case <-time.After(500 * time.Millisecond):
	}
}