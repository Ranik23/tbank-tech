package kafkaproducer

import (
	"encoding/json"
	"log/slog"

	"github.com/IBM/sarama"
	"github.com/google/go-github/v69/github"
)

type KafkaProducer struct {
	logger			*slog.Logger
	sarama.AsyncProducer
	commitCh 		chan *github.RepositoryCommit
	stopCh			chan struct{}
}

func (kp *KafkaProducer) Run() {
	go func() {
		for {
			select {
			case commit, ok := <-kp.commitCh:
				if !ok {
					// Канал commitCh закрыт, завершаем работу
					return
				}
				topic := "test"
				kp.produceCommit(commit, topic)
				//
			case <- kp.stopCh:
				return
			}
		}
	}()
}


func (kp *KafkaProducer) Stop() {
	kp.logger.Info("Kafka Producer stopped")
	kp.stopCh <- struct{}{}
}


func (kp *KafkaProducer) produceCommit(commit *github.RepositoryCommit, topic string) error {
	commitJSON, err := json.Marshal(*commit)
	if err != nil {
		kp.logger.Error("Failed to marshall the commit: %v", err)
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(commitJSON),
	}

	kp.Input() <- msg

	select {
	case err := <-kp.Errors():

	case succeses := <-kp.Successes():
		
	}


	return nil
}