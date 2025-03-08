package clientkafka_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"tbank/scrapper/internal/hub"
	"tbank/scrapper/pkg/github/utils"
	"testing"
	"time"

	git "tbank/scrapper/internal/mocks/github"

	"github.com/IBM/sarama"
	"github.com/golang/mock/gomock"
	"github.com/google/go-github/v69/github"
	"github.com/stretchr/testify/assert"
)

const (
	InfoEmoji        = "\u2139"  // ℹ️ Information
	CheckMarkEmoji   = "\u2705"  // ✅ Check Mark
	CrossMarkEmoji   = "\u274C"  // ❌ Cross Mark
	RocketEmoji      = "\u1F680" // 🚀 Rocket
	FireEmoji        = "\u1F525" // 🔥 Fire
	HighVoltageEmoji = "\u26A1"  // ⚡ High Voltage
	HammerEmoji      = "\u1F6E0" // 🛠️ Hammer and Wrench
	StarEmoji        = "\u2B50"  // 🌟 Glowing Star
)

func TestClientKafka_PublishCommit(t *testing.T) {
	t.Logf("%s Starting TestClientKafka_PublishCommit", InfoEmoji)

	addresses := []string{"0.0.0.0:9093"}

	url := "https://github.com/MAtveyka12/tbank"

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Consumer.Return.Errors = true

	asyncProducer, err := sarama.NewAsyncProducer(addresses, config)
	if !assert.NoErrorf(t, err, "%s Failed to create a new AsyncProducer: %v", CrossMarkEmoji, err) {
		return
	}
	defer asyncProducer.Close()

	consumer, err := sarama.NewConsumer(addresses, config)
	if !assert.NoError(t, err, "%s Failed to create a Kafka Consumer: %v", CrossMarkEmoji, err) {
		return
	}
	defer consumer.Close()

	t.Logf("%s Adding URL: %s", FireEmoji, url)

	owner, repo, err := utils.GetLinkParams(url)
	if err != nil {
		t.Fatalf("%s Wrong URL scheme: %v", CrossMarkEmoji, err)
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockGitHubClient := git.NewMockGitHubClient(ctrl)
	mockGitHubClient.EXPECT().LatestCommit(gomock.Any(), owner, repo, gomock.Any()).Return(
		&github.RepositoryCommit{
			SHA: github.Ptr("test_sha"),
			Commit: &github.Commit{
				Message: github.Ptr("Test commit message"),
			},
		}, nil, nil).AnyTimes()

	client := hub.NewClient(asyncProducer, slog.Default(), url, mockGitHubClient)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	go func(t *testing.T) {
		if err := client.Run(ctx, 10 * time.Second); err != nil {
			t.Fatalf("%s Failed to run the client: %v", CrossMarkEmoji, err)
		}
	}(t)

	topicName := fmt.Sprintf("%s_%s", owner, repo)

	partitionConsumer, err := consumer.ConsumePartition(topicName, 0, sarama.OffsetNewest)
	if err != nil {
		t.Fatalf("%s Failed to consume partition: %v", CrossMarkEmoji, err)
	}
	defer partitionConsumer.Close()

	message := <-partitionConsumer.Messages()

	var msg map[string]interface{}
	if err := json.Unmarshal(message.Value, &msg); err != nil {
		t.Fatalf("%s Failed to unmarshall data: %v", CrossMarkEmoji, err)
	}

	assert.Equalf(t, "Test commit message", msg["message"], "%s Unexpected commit message", CrossMarkEmoji)
	assert.Equalf(t, "test_sha", msg["sha"], "%s Unexpected commit SHA", CrossMarkEmoji)


	t.Logf("%s TestClientKafka_PublishCommit Passed", CheckMarkEmoji)
}
