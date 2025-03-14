package kafkaconsumer

import (
	"log/slog"
	"github.com/IBM/sarama"
)

type KafkaConsumer struct {
	consumer sarama.Consumer
	stopCh	 chan struct{}
	logger	 *slog.Logger
}

func NewKafkaConsumer(kafkaConsumer sarama.Consumer, messageCh chan sarama.ConsumerMessage) (*KafkaConsumer, error) {
	return &KafkaConsumer{
		consumer: kafkaConsumer,
	}, nil
}

func (kc *KafkaConsumer) Run() {
	go func() {
		for {

		}
	}()
}

func (kc *KafkaConsumer) Stop() {
	close(kc.stopCh)
}