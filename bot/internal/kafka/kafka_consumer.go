package kafkaconsumer

import (
	"log/slog"

	"github.com/IBM/sarama"
)

type KafkaConsumer struct {
	consumer sarama.Consumer
	commitCh chan sarama.ConsumerMessage
	stopCh   chan struct{}
	logger   *slog.Logger
	topic    string
}

func NewKafkaConsumer(addresses []string, topicToRead string,
						commitCh chan sarama.ConsumerMessage, logger *slog.Logger, saramaConfig *sarama.Config) (*KafkaConsumer, error) {
	consumer, err := sarama.NewConsumer(addresses, saramaConfig)
	if err != nil {
		logger.Error("Failed to create a new Sarama consumer", slog.String("error", err.Error()))
		return nil, err
	}
	return &KafkaConsumer{
		consumer: consumer,
		logger:   logger,
		topic:    topicToRead,
		commitCh: commitCh,
		stopCh:   make(chan struct{}),
	}, nil
}

func (kc *KafkaConsumer) Run() error {
	const op = "KafkaConsumer.Run"

	kc.logger.Info(op, slog.String("message", "Kafka Consumer Running"))

	partitions, err := kc.consumer.Partitions(kc.topic)
	if err != nil {
		kc.logger.Error(op, slog.Any("error", err))
		return err
	}

	for _, partition := range partitions {
		go kc.consumePartition(partition)
	}
	return nil
}

func (kc *KafkaConsumer) consumePartition(partition int32) {
	const op = "KafkaConsumer.consumePartition"

	partitionConsumer, err := kc.consumer.ConsumePartition(kc.topic, partition, sarama.OffsetNewest)
	if err != nil {
		kc.logger.Error(op, slog.Any("error", err))
		return
	}
	defer partitionConsumer.Close()

	kc.handleMessages(partitionConsumer)
}

func (kc *KafkaConsumer) handleMessages(consumer sarama.PartitionConsumer) {
	const op = "KafkaConsumer.handleMessages"

	for {
		select {
		case msg, ok := <-consumer.Messages():
			if !ok {
				kc.logger.Warn(op, slog.String("message", "partitionConsumer.Messages() closed"))
				return
			}
			kc.commitCh <- *msg
			kc.logger.Info(op, slog.String("message", "Message sent to commit channel"))

		case err, ok := <-consumer.Errors():
			if ok {
				kc.handleErrors(err)
			}

		case <-kc.stopCh:
			kc.logger.Info(op, slog.String("message", "Kafka Consumer stopped"))
			return
		}
	}
}

func (kc *KafkaConsumer) handleErrors(err *sarama.ConsumerError) {
	const op = "KafkaConsumer.handleErrors"
	kc.logger.Error(op, slog.Any("error", err))
}

func (kc *KafkaConsumer) Stop() {
	close(kc.stopCh)
	if err := kc.consumer.Close(); err != nil {
		kc.logger.Error("KafkaConsumer.Stop", slog.Any("error", err))
	}
}
