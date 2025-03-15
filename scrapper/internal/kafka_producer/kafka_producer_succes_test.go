package kafkaproducer

import (
	"log/slog"
	"sync"
	"tbank/scrapper/internal/hub"
	"testing"
	"github.com/IBM/sarama/mocks"
	"github.com/google/go-github/v69/github"
)




func TestKafkaProducerSucces(t *testing.T) {
    commitCh := make(chan hub.CustomCommit)
    logger := slog.Default()

    mockProducer := mocks.NewAsyncProducer(t, nil)

    var wg sync.WaitGroup
	wg.Add(1)

	checkMessage := func(msg []byte) error {
		wg.Done()
		return nil
	}

	mockProducer.ExpectInputWithCheckerFunctionAndSucceed(checkMessage)
	

    kafkaProducer := &KafkaProducer{
        logger:   logger,
        producer: mockProducer,
        commitCh: commitCh,
        topicToSend: "test",
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


	wg.Wait()


    kafkaProducer.Stop()
}
