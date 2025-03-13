package kafkaproducer

import (
	"encoding/json"
	"log/slog"
	"strconv"
	"tbank/scrapper/internal/hub"

	"github.com/IBM/sarama"
)

type KafkaProducer struct {
	logger   *slog.Logger
	producer sarama.AsyncProducer
	commitCh chan hub.CustomCommit
	stopCh   chan struct{}
}

func NewKafkaProducer(producer sarama.AsyncProducer, logger *slog.Logger, commitCh chan hub.CustomCommit) (*KafkaProducer, error) {
	return &KafkaProducer{
		logger:   logger,
		producer: producer,
		commitCh: commitCh,
		stopCh:   make(chan struct{}),
	}, nil
}

func (kp *KafkaProducer) Run() {
	const op = "KafkaProducer.Run"
	kp.logger.Info(op, slog.String("msg", "Kafka producer is running"))
	go func() {
		for {
			select {
			case commit, ok := <-kp.commitCh:
				if !ok {
					kp.logger.Warn(op, slog.String("msg", "Commit channel closed, stopping producer"))
					return
				}
				topicUserID := strconv.Itoa(int(commit.UserID))

				if err := kp.produceCommit(commit, topicUserID); err != nil {
					kp.logger.Error(op, slog.String("topic", topicUserID), slog.String("error", err.Error()))
				} else {
					kp.logger.Info(op, slog.String("user_id", topicUserID), slog.String("msg", "Produced commit successfully"))
				}

			case <-kp.stopCh:
				kp.logger.Warn(op, slog.String("msg", "Stopping Kafka producer"))
				return
			}
		}
	}()
}

func (kp *KafkaProducer) Stop() {
	const op = "KafkaProducer.Stop"
	kp.logger.Warn(op, slog.String("msg", "Stopping Kafka producer"))
	kp.producer.Close()
	close(kp.stopCh)
	close(kp.commitCh)
}

func (kp *KafkaProducer) produceCommit(commit hub.CustomCommit, topic string) error {
	const op = "KafkaProducer.produceCommit"
	commitJSON, err := json.Marshal(commit)
	if err != nil {
		kp.logger.Error(op, slog.String("error", err.Error()))
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(commitJSON),
	}

	kp.producer.Input() <- msg

	select {
	case errMsg := <-kp.producer.Errors():
		kp.logger.Error(op, slog.String("error", errMsg.Err.Error()))
		return errMsg.Err
	case <-kp.producer.Successes():
		kp.logger.Info(op, slog.String("msg", "Successfully sent the message"))
		return nil
	}
}
