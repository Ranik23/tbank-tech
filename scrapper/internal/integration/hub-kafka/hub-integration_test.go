package hub_kafka_test

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"tbank/scrapper/internal/hub"
	git "tbank/scrapper/internal/mocks/github"
	"tbank/scrapper/pkg/github/utils"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/golang/mock/gomock"
	"github.com/google/go-github/v69/github"
	"github.com/stretchr/testify/assert"
)

func TestHub_CommitsSuccess(t *testing.T) {
	t.Log("ℹ️ Starting TestHub_CommitsSuccess")

	addresses := []string{"0.0.0.0:9093"}

	urls := []string{
		"https://github.com/MAtveyka12/tbank",
		"https://github.com/v1lezz/WarehouseLamoda",
		"https://github.com/epchamp001/AVIP-MEPhI-2025",
	}
	
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Consumer.Return.Errors = true

	asyncProducer, err := sarama.NewAsyncProducer(addresses, config)
	if !assert.NoError(t, err, "❌ Failed to create a new AsyncProducer") {
		return
	}
	defer asyncProducer.Close()

	consumer, err := sarama.NewConsumer(addresses, config)
	if !assert.NoError(t, err, "❌ Failed to create a Kafka Consumer") {
		return
	}
	defer consumer.Close()

	for _, url := range urls {

		t.Logf("🔹 Adding URL: %s", url)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		owner, repo, err := utils.GetLinkParams(url)
		if err != nil {
			t.Fatalf("❌ Wrong URL scheme")
		}
		t.Logf("🔹 Owner - %s, Repo - %s", owner, repo)

		topicName := fmt.Sprintf("%s/%s", owner, repo)

		t.Logf("🔹 Topic name - %s", topicName)

		mockGitHubClient := git.NewMockGitHubClient(ctrl)
		mockGitHubClient.EXPECT().LatestCommit(gomock.Any(), owner, repo, gomock.Any()).Return(
			&github.RepositoryCommit{
					SHA: github.Ptr("test_sha"),
					Commit: &github.Commit{
						Message: github.Ptr("Test commit message"),
					},
				}, nil, nil).Times(1)

				
		hub := hub.NewHub(asyncProducer, slog.Default(), mockGitHubClient)

		if err := hub.AddTrack(url, 1); err != nil {
			t.Fatalf("Failed to add the track: %v", err)
		}

		partitionConsumer, err := consumer.ConsumePartition(topicName, 0, sarama.OffsetNewest)
		if !assert.NoError(t, err, "❌ Failed to consume partition") {
			return
		}
		defer partitionConsumer.Close()

		var message *sarama.ConsumerMessage
		select {
		case msg := <-partitionConsumer.Messages():
			message = msg
		case <-time.After(5 * time.Second):
			t.Fatalf("❌ Timeout: No message received from Kafka for topic %s", url)
		}

		var msgData map[string]interface{}
		if err := json.Unmarshal(message.Value, &msgData); err != nil {
			t.Fatalf("❌ Failed to unmarshal message: %v", err)
		}

		assert.Equal(t, "Test commit message", msgData["message"], "❌ Unexpected commit message")
		assert.Equal(t, "test_sha", msgData["sha"], "❌ Unexpected commit SHA")
		assert.Equal(t, url, msgData["url"], "❌ Unexpected URL")

		t.Logf("✅ Test passed for URL: %s", url)
	}
}
