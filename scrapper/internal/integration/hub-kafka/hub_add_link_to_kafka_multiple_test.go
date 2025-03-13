package hubkafka

import (
	"encoding/json"
	"log/slog"
	"tbank/scrapper/internal/hub"
	kafkaproducer "tbank/scrapper/internal/kafka-producer"
	git "tbank/scrapper/internal/mocks/github"
	"testing"
	"github.com/IBM/sarama"
	"github.com/golang/mock/gomock"
	"github.com/google/go-github/v69/github"
)

func TestHub_AddLink_SendsCommitToKafkaAndReceivesItMultiple(t *testing.T) {

//											SETTING UP AND STARTING HUB AND KAFKA CONSUMER											//
//----------------------------------------------------------------------------------------------------------------------------------//

	addresses := []string{"localhost:9093"}
	linkExample := "https://github.com/Ranik23/tbank-tech"

	commitExample1 := &github.RepositoryCommit{
		SHA: github.Ptr("test_sha"),
		Commit: &github.Commit{
			Message: github.Ptr("test_message"),
		},
	}

	commitExample2 := &github.RepositoryCommit{
		SHA: github.Ptr("test_sha2"),
		Commit: &github.Commit{
			Message: github.Ptr("test_message2"),
		},
	}

	commitExample3 := &github.RepositoryCommit{
		SHA: github.Ptr("test_sha3"),
		Commit: &github.Commit{
			Message: github.Ptr("test_message3"),
		},
	}

	exampleCommits := []*github.RepositoryCommit{
		commitExample1,
		commitExample2,
		commitExample3,
	}

	logger := slog.Default()

	controller := gomock.NewController(t)

	mockGit := git.NewMockGitHubClient(controller)

	call1 := mockGit.EXPECT().LatestCommit(gomock.Any(), "Ranik23", "tbank-tech", gomock.Any()).Return(
		commitExample1, nil, nil,
	).Times(1)

	call2 := mockGit.EXPECT().LatestCommit(gomock.Any(), "Ranik23", "tbank-tech", gomock.Any()).Return(
		commitExample2, nil, nil,
	).Times(1).After(call1)

	mockGit.EXPECT().LatestCommit(gomock.Any(), "Ranik23", "tbank-tech", gomock.Any()).Return(
		commitExample3, nil, nil,
	).AnyTimes().After(call2)


	commitCh := make(chan hub.CustomCommit)

	myHub := hub.NewHub(mockGit, commitCh, logger)

	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.Return.Errors = true

	producer, err := sarama.NewAsyncProducer(addresses, saramaConfig)
	if err != nil {
		t.Fatalf("Failed to create a new async producer: %v", err)
	}

	kafkaProducer, err := kafkaproducer.NewKafkaProducer(producer, logger, commitCh)
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

	partitions, err := kafkaConsumer.Partitions("1")
	if err != nil {
		t.Fatalf("Failed to list the partitions for the given topic: %v", err)
	}

	if len(partitions) == 0 {
		t.Fatalf("Partitions number should be > 0")
	}

	t.Logf("Partitions - %d", len(partitions))

	theOnlyPartition := partitions[0]

	partitionConsumer, err := kafkaConsumer.ConsumePartition("1", theOnlyPartition, sarama.OffsetNewest)
	if err != nil {
		t.Fatalf("Failed to create a consumer on this specific partition: %v", err)
	}
	defer partitionConsumer.Close()

	errorCh := make(chan error)
	receivedCommitCh := make(chan hub.CustomCommit)

	go func() {
		for message := range partitionConsumer.Messages() {
			var receivedCommit hub.CustomCommit
			if err := json.Unmarshal(message.Value, &receivedCommit); err != nil {
				errorCh <- err
				return
			}
			logger.Info("Message sent to receivedChannel")
			receivedCommitCh <- receivedCommit
		}
	}()

	i := 0

Loop:
	for i <= 2 {
		select {
		case err := <-errorCh:
			t.Fatalf("Error occurred: %v", err)
		case commit := <-receivedCommitCh:
			if *commit.Commit.SHA != *exampleCommits[i].SHA ||
				*commit.Commit.Commit.Message != *exampleCommits[i].Commit.Message {
				t.Fatalf("Commit Message is not the same. Expected: %s %s, Got: %s %s",
					*exampleCommits[i].SHA, *exampleCommits[i].Commit.Message,
					*commit.Commit.SHA, *commit.Commit.Commit.Message)
			}

			t.Logf("Expected SHA - %s, Real SHA - %s, Expected Message - %s, Real Message - %s",
				*exampleCommits[i].SHA, *commit.Commit.SHA,
				*exampleCommits[i].Commit.Message, *commit.Commit.Commit.Message)

			if i == 2 {
				break Loop
			}

			i++
		}
	}

	close(errorCh)
	close(receivedCommitCh)

}
