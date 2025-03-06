package hub_kafka_test

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"tbank/scrapper/internal/hub"
	git "tbank/scrapper/internal/mocks/github"
	"testing"

	"github.com/IBM/sarama"
	"github.com/golang/mock/gomock"
	"github.com/google/go-github/v69/github"
	"github.com/stretchr/testify/assert"
)

func TestHub_CommitsSucces(t *testing.T) {

	addresses := []string{"0.0.0.0:9093"}

	asyncProducer, err := sarama.NewAsyncProducer(addresses, nil)
	if !assert.NoError(t, err) {
		t.Fatalf("error: %v", err)
	}
	defer asyncProducer.Close() 




	ctrl := gomock.NewController(t)
	mockGitHubClient := git.NewMockGitHubClient(ctrl)
	mockGitHubClient.EXPECT().ListCommits(gomock.Any(), "test_owner", "test_repo", gomock.Any()).Return(
		[]*github.RepositoryCommit{
			{
				SHA: github.Ptr("test_sha"),
				Commit: &github.Commit{
					Message: github.Ptr("Test commit message"),
				},
			},
		}, nil, nil).AnyTimes()




	hub := hub.NewHub(asyncProducer, slog.Default(), mockGitHubClient, "test_topic")
	hub.AddTrack("https://github.com/test_owner/test_repo")




	consumer, err := sarama.NewConsumer(addresses, nil)
	if !assert.NoError(t, err) {
		t.Fatalf("Failed to create a New Consumer: %v", err)
	}
	defer consumer.Close()





	partitions, err := consumer.Partitions("test_topic")
	if !assert.NoError(t, err) {
		t.Fatalf("Failed to get the partitions: %v", err)
	}


	

	var wg sync.WaitGroup
	errorCh := make(chan error, 1)  
	doneCh := make(chan struct{})  
	var once sync.Once   
	var messageReceived bool = false 

	go func() {
		select {
		case err := <-errorCh:
			if err != nil {
				t.Errorf("Test failed due to error: %v", err)
			}
		case <-doneCh:
		}
	}()

	for _, partition := range partitions {
		consPartition, err := consumer.ConsumePartition("test_topic", partition, sarama.OffsetNewest)
		if !assert.NoError(t, err) {
			t.Fatalf("Failed to consume partition: %v", err)
		}

		wg.Add(1)
		go func(consumer sarama.PartitionConsumer) {
			defer wg.Done()
			defer consPartition.Close()

			for {
				select {
				case message := <-consumer.Messages():
					var msg map[string]interface{}
					if err := json.Unmarshal(message.Value, &msg); err != nil {
						once.Do(func() { errorCh <- fmt.Errorf("Failed to unmarshal message: %v", err) })
						return
					}

					fmt.Println(msg)

					commitMsg, ok1 := msg["message"].(string)
					sha, ok2 := msg["sha"].(string)

					if !ok1 || !ok2 {
						once.Do(func() { errorCh <- fmt.Errorf("Invalid message format: %v", msg) })
						return
					}

					if commitMsg == "Test commit message" && sha == "test_sha" {
						messageReceived = true
						once.Do(func() { close(doneCh) })
						return
					}
				case <-doneCh:
					return
				}
			}
		}(consPartition)
	}
	wg.Wait()

	if !messageReceived {
		t.Fatal("Expected message not received in any partition")
	}
}
