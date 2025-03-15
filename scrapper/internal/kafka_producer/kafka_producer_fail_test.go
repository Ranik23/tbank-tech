package kafkaproducer

import (
	"log/slog"
	"tbank/scrapper/internal/hub"
	"testing"


	"github.com/IBM/sarama"
	"github.com/IBM/sarama/mocks"
	"github.com/google/go-github/v69/github"
)



func TestKafkaProducer_ProduceCommitFailure(t *testing.T) {

	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.Return.Errors = true


    commitCh := make(chan hub.CustomCommit, 1)
    logger := slog.Default()

    mockProducer := mocks.NewAsyncProducer(t, saramaConfig)

    mockProducer.ExpectInputAndFail(sarama.ErrOutOfBrokers)

    kafkaProducer := &KafkaProducer{
        logger:   logger,
        producer: mockProducer,
        commitCh: commitCh,
		stopCh: make(chan struct{}),
    }

    kafkaProducer.Run()

    testCommit := hub.CustomCommit{
        UserID: 123,
        Commit: &github.RepositoryCommit{
            SHA: github.Ptr("testsha"),
            Commit: &github.Commit{
                Message: github.Ptr("Test commit"),
            },
        },
    }

    commitCh <- testCommit

    select {
    case errMsg := <-mockProducer.Errors():
        if errMsg.Err != sarama.ErrOutOfBrokers {
            t.Fatalf("Expected error %v but got %v", sarama.ErrOutOfBrokers, errMsg.Err)
        }
    case <-mockProducer.Successes():
        t.Fatal("Expected message to fail, but it succeeded")
    }

    kafkaProducer.Stop()
}
