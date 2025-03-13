package kafkaconsumer

import (
	"context"
	"log/slog"
	"time"

	"github.com/IBM/sarama"
)

type KafkaConsumer struct {
	logger   *slog.Logger
	consumer sarama.Consumer
	stopCh   chan struct{}
	topics   []string
	client   sarama.Client
}

func NewKafkaConsumer(brokers []string, logger *slog.Logger) (*KafkaConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	client, err := sarama.NewClient(brokers, config)
	if err != nil {
		return nil, err
	}

	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		client.Close()
		return nil, err
	}

	return &KafkaConsumer{
		logger:   logger,
		consumer: consumer,
		stopCh:   make(chan struct{}),
		client:   client,
	}, nil
}

func (k *KafkaConsumer) updateTopicsPeriodically(ctx context.Context, interval time.Duration) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(interval):
				topics, err := k.client.Topics()
				if err != nil {
					k.logger.Error("Ошибка получения списка топиков", slog.String("error", err.Error()))
					continue
				}

				k.topics = topics
				k.logger.Info("Обновлён список топиков", slog.Any("topics", topics))
			}
		}
	}()
}

func (k *KafkaConsumer) Run(ctx context.Context) error {

	go k.updateTopicsPeriodically(ctx, 10*time.Second)

	errorCh := make(chan error)

	go func() {

	for {
		select {
		case <-k.stopCh:
			k.logger.Info("Kafka Consumer остановлен")
			return nil
		default:
			for _, topic := range k.topics {
				partitions, err := k.consumer.Partitions(topic)
				if err != nil {
					k.logger.Error("Ошибка получения партиций", slog.String("error", err.Error()))
					continue
				}

				for _, partition := range partitions {
					go k.consumePartition(topic, partition, errorCh)
				}
			}
		}

		select {
		case err := <-errorCh:
			k.logger.Error("Ошибка во время потребления сообщений", slog.String("error", err.Error()))
		case <-ctx.Done():
			k.logger.Info("Kafka Consumer: контекст завершён")
			return nil
		}

		time.Sleep(5 * time.Second)
	}
	}()
}

func (k *KafkaConsumer) consumePartition(topic string, partition int32, errorCh chan error) {
	partitionConsumer, err := k.consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
	if err != nil {
		errorCh <- err
		return
	}
	defer partitionConsumer.Close()

	for msg := range partitionConsumer.Messages() {
		k.logger.Info("Получено сообщение",
			slog.String("topic", msg.Topic),
			slog.Int64("offset", msg.Offset),
			slog.String("value", string(msg.Value)))
	}
}

func (k *KafkaConsumer) Stop() {
	close(k.stopCh)
	k.consumer.Close()
	k.client.Close()
}
