package hubkafka





import (
	"encoding/json"
	"log/slog"
	"tbank/scrapper/internal/hub"
	kafkaproducer "tbank/scrapper/internal/kafka_producer"
	git "tbank/scrapper/pkg/github_client/mock"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/golang/mock/gomock"
	"github.com/google/go-github/v69/github"
)

func TestHub_AddLink_SendsCommitToKafkaAndReceivesIt(t *testing.T) {

//											SETTING UP AND STARTING HUB AND KAFKA CONSUMER											//
//----------------------------------------------------------------------------------------------------------------------------------//
	
	addresses := []string{"localhost:9093"}
	linkExample := "https://github.com/Ranik23/tbank-tech"
	exampleTopic := "test_topic"

	commitExample := &github.RepositoryCommit{
		SHA: github.Ptr("test_sha"),
		Commit: &github.Commit{
			Message: github.Ptr("test_message"),
		},
	}

	logger := slog.Default()

	controller := gomock.NewController(t)

	mockGit := git.NewMockGitHubClient(controller)

	mockGit.EXPECT().LatestCommit(gomock.Any(), "Ranik23", "tbank-tech", gomock.Any()).Return(
		commitExample, nil, nil,
	).AnyTimes()

	commitCh := make(chan hub.CustomCommit)

	myHub := hub.NewHub(mockGit, commitCh, logger)

	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.Return.Errors = true

	producer, err := sarama.NewAsyncProducer(addresses, saramaConfig)
	if err != nil {
		t.Fatalf("Failed to create a new async producer: %v", err)
	}

	kafkaProducer, err := kafkaproducer.NewKafkaProducer(producer, logger, commitCh, exampleTopic)
	if err != nil {
		t.Fatalf("Failed to create a new kafka: %v", err)
	}

	kafkaProducer.Run()
	defer kafkaProducer.Stop()

	myHub.Run()
	defer myHub.Stop()

	myHub.AddLink(linkExample, 1)

//															CONSUMER STARTING                                                       // 
//----------------------------------------------------------------------------------------------------------------------------------//

	kafkaConsumer, err := sarama.NewConsumer(addresses, saramaConfig)
	if err != nil {
		t.Fatalf("Failed to create a new consumer: %v", err)
	}
	defer kafkaConsumer.Close()

	partitions, err := kafkaConsumer.Partitions(exampleTopic)
	if err != nil {
		t.Fatalf("Failed to list the partitions for the given topic: %v", err)
	}

	stopCh := make(chan struct{}, len(partitions))
	errorCh := make(chan error, len(partitions))

	t.Logf("Partitions - %d", len(partitions))
	receivedCommitCh := make(chan hub.CustomCommit, 1)

	for _, partition := range partitions {
		go func(partition int32) {
			partitionConsumer, err := kafkaConsumer.ConsumePartition(exampleTopic, partition, sarama.OffsetNewest)
			if err != nil {
				errorCh <- err
				return
			}
			defer partitionConsumer.Close()

			for {
				select {
				case msg := <-partitionConsumer.Messages():
					var receivedCommit hub.CustomCommit
					if err := json.Unmarshal(msg.Value, &receivedCommit); err != nil {
						errorCh <- err
						return
					}

					receivedCommitCh <- receivedCommit

					stopCh <- struct{}{}

					return

				case <-stopCh:
					return
				}
			}
		}(partition)
	}

	select {
	case err := <-errorCh:
		t.Fatalf("Error occurred: %v", err)
	case commit := <-receivedCommitCh:
		if *commit.Commit.SHA != *commitExample.SHA || *commit.Commit.Commit.Message != *commitExample.Commit.Message {
			t.Fatalf("Commit Message is not the same. What I Got - %s %s", *commit.Commit.SHA, *commit.Commit.Commit.Message)
		}
		t.Logf("Expected SHA - %s, Real SHA - %s, Expected Message - %s, Real Message - %s", *commit.Commit.SHA, *commitExample.SHA,
			*commit.Commit.Commit.Message, *commitExample.Commit.Message)
	case <-time.After(5 * time.Second):
		t.Fatalf("Timeout expired")
	}

	close(stopCh)
	close(receivedCommitCh)
	close(errorCh)
}

