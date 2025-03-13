package kafkaproducer

import (
	"log/slog"
	"tbank/scrapper/internal/hub"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/IBM/sarama/mocks"
	"github.com/google/go-github/v69/github"
)


func TestKafkaProducer_Stop(t *testing.T) {
    saramaConfig := sarama.NewConfig()
    saramaConfig.Producer.Return.Successes = true
    saramaConfig.Producer.Return.Errors = true


    commitCh := make(chan hub.CustomCommit, 1)
    logger := slog.Default()

    mockProducer := mocks.NewAsyncProducer(t, saramaConfig)

    mockProducer.ExpectInputAndSucceed()

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
    case <-mockProducer.Successes():
        t.Log("Producer message succeeded")
    case <-mockProducer.Errors():
        t.Fatal("Expected message to succeed, but it failed")
    case <-time.After(5 * time.Second):
        t.Fatal("Timed out waiting for success or error")
    }

    kafkaProducer.Stop()

     _, ok := <-mockProducer.Successes()
    if ok {
        t.Fatal("Expected Successes channel to be closed, but it was open")
    } else {
        t.Log("Successes channel is closed after Stop")
    }

    _, ok = <-mockProducer.Errors()
    if ok {
        t.Fatal("Expected Errors channel to be closed, but it was open")
    } else {
        t.Log("Errors channel is closed after Stop")
    }
}
