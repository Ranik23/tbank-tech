package kafkaconsumer

import (
	"errors"
	"testing"
	"time"
	"github.com/IBM/sarama"
	"github.com/IBM/sarama/mocks"
	"github.com/stretchr/testify/require"
)





func TestKafkaConsumer_ErrorOnMessage(t *testing.T) {
	exampleTopic := "test"
	examplePartitions := []int32{0}

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true  // Включаем ошибки

	mockConsumer := mocks.NewConsumer(t, config)
	mockConsumer.SetTopicMetadata(map[string][]int32{exampleTopic: examplePartitions})

	mockPartitionConsumer := mockConsumer.ExpectConsumePartition(exampleTopic, examplePartitions[0], sarama.OffsetNewest)

	// Отправляем одно валидное сообщение
	mockPartitionConsumer.YieldMessage(&sarama.ConsumerMessage{Topic: exampleTopic, Value: []byte("valid")})

	// Отправляем ошибку
	mockPartitionConsumer.YieldError(errors.New("simulated error"))

	commitCh := make(chan sarama.ConsumerMessage, 1)
	kafkaConsumer := NewKafkaConsumer(mockConsumer, exampleTopic, commitCh)


	go kafkaConsumer.Run()
	defer kafkaConsumer.Stop()

	receivedMsg := false
	timeout := time.After(2 * time.Second)

	for {
		select {
		case msg := <-commitCh:
			require.Equal(t, "valid", string(msg.Value))
			receivedMsg = true
		case <-timeout:
			if !receivedMsg {
				t.Fatal("Timeout waiting for valid message")
			}
			return
		}
	}
}
