package kafkacosumer

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)


const GroupID = "my-single-consumer-group"

type KafkaConsumer struct {
	consumer *kafka.Consumer
	topics   []string
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewKafkaConsumer(brokers, topics []string) (*KafkaConsumer, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": 	brokers,
		"group.id":         	GroupID ,
		"auto.offset.reset": 	"latest",
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка создания consumer: %w", err)
	}

	err = consumer.SubscribeTopics(topics, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка подписки: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &KafkaConsumer{
		consumer: consumer,
		topics:   topics,
		ctx:      ctx,
		cancel:   cancel,
	}, nil
}

func (kc *KafkaConsumer) Run() {
	go func() {
		for {
			select {
			case <-kc.ctx.Done():
				log.Println("Kafka Consumer: остановка по сигналу")
				return
			default:
				msg, err := kc.consumer.ReadMessage(1 * time.Second)
				if err == nil {
					fmt.Printf("Получено сообщение: %s -> %s\n", msg.TopicPartition, string(msg.Value))
				} else if err.(kafka.Error).Code() != kafka.ErrTimedOut {
					log.Printf("Ошибка чтения сообщения: %v\n", err)
				}
			}
		}
	}()
}

func (kc *KafkaConsumer) Stop() {
	log.Println("Остановка Kafka Consumer...")
	kc.cancel()
	kc.consumer.Close()
	log.Println("Kafka Consumer остановлен")
}