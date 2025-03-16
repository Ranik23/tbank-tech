package kafkaconsumer

import (
	"log/slog"
	"github.com/IBM/sarama"
)

type KafkaConsumer struct {
	consumer sarama.Consumer
	commitCh chan sarama.ConsumerMessage
	stopCh	 chan struct{}
	logger	 *slog.Logger
	topic	 string
}

func NewKafkaConsumer(kafkaConsumer sarama.Consumer, topicToRead string, commitCh chan sarama.ConsumerMessage, logger *slog.Logger) *KafkaConsumer {
	return &KafkaConsumer{
		consumer: kafkaConsumer,
		logger: logger,
		topic: topicToRead,
		commitCh: commitCh,
		stopCh: make(chan struct{}),
	}
}

func (kc *KafkaConsumer) Run() {
	const op = "KafkaConsumer.Run"

	partitions, err := kc.consumer.Partitions(kc.topic)
	if err != nil {
		kc.logger.Error(op, slog.String("message", err.Error()))	
		return
	}

	for _, partition := range partitions {
		go func() {
			partitionConsumer, err := kc.consumer.ConsumePartition(kc.topic, partition, sarama.OffsetNewest)
			if err != nil {
				kc.logger.Error(op, slog.String("message", err.Error()))
				return
			}

			defer partitionConsumer.Close()

			for {
				select {
				case msg, ok := <-partitionConsumer.Messages():
					if !ok {
						kc.logger.Warn(op, slog.String("message", "partitionConsumer.Messages() closed"))
						return
					}
					kc.commitCh <- *msg
					kc.logger.Info(op, slog.String("message", "Message sent to commit channel"))

				case err, ok := <-partitionConsumer.Errors():
					if ok {
						kc.logger.Error(op, slog.String("error", err.Error()))
					}
				case <-kc.stopCh:
					kc.logger.Info(op, slog.String("message", "Kafka Consumer stopped"))
					return
				}
			}
		}()
	}
}

func (kc *KafkaConsumer) Stop() {
	close(kc.stopCh)
	kc.consumer.Close()
}