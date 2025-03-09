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
	stopCH 			chan struct{}
}

func (kp *KafkaProducer) Run() {
	go func() {
		for {
			select {
			case _, ok := <-kp.commitCh:
				if !ok {
					// Канал commitCh закрыт, завершаем работу
					return
				}
				//
			case <- kp.stopCH:
				return
			}
		}
	}()
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