package hub_kafka_test

import (
	"encoding/json"
	"log/slog"
	"tbank/scrapper/internal/hub"
	git "tbank/scrapper/internal/mocks/github"
	"testing"

	"github.com/IBM/sarama"
	"github.com/golang/mock/gomock"
	"github.com/google/go-github/v69/github"
	"github.com/stretchr/testify/assert"
)

func TestHub_CommitsSuccess(t *testing.T) {
	addresses := []string{"0.0.0.0:9093"}

	asyncProducer, err := sarama.NewAsyncProducer(addresses, nil)
	if !assert.NoError(t, err) {
		t.Fatalf("❌ Failed to create a new AsyncProducer: %v", err)
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
		t.Fatalf("❌ Failed to create a New Consumer: %v", err)
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition("test_topic", 0, sarama.OffsetNewest)
	if !assert.NoError(t, err) {
		t.Fatalf("❌ Failed to consume partition: %v", err)
	}
	defer partitionConsumer.Close()

	message := <-partitionConsumer.Messages()

	var msg map[string]interface{}
	if err := json.Unmarshal(message.Value, &msg); err != nil {
		t.Fatalf("❌ Failed to unmarshal message: %v", err)
	}

	commitMsg, ok1 := msg["message"].(string)
	sha, ok2 := msg["sha"].(string)

	if !ok1 || !ok2 {
		t.Fatalf("❌ Invalid message format: %v", msg)
	}

	if commitMsg != "Test commit message" || sha != "test_sha" {
		t.Fatalf("❌ Unexpected message content: %+v", msg)
	}

	t.Log("✅ Test passed: Expected message received")
}

func TestHub_CommitsFail(t *testing.T) {
	addresses := []string{"0.0.0.0:9093"}

	asyncProducer, err := sarama.NewAsyncProducer(addresses, nil)
	if !assert.NoError(t, err) {
		t.Fatalf("❌ Failed to create a new AsyncProducer: %v", err)
	}
	defer asyncProducer.Close()

	ctrl := gomock.NewController(t)
	mockGitHubClient := git.NewMockGitHubClient(ctrl)
	mockGitHubClient.EXPECT().ListCommits(gomock.Any(), "test_owner", "test_repo", gomock.Any()).Return(
		[]*github.RepositoryCommit{
			{
				SHA: github.Ptr("unexpected_sha"), // ⚠️ Не тот SHA
				Commit: &github.Commit{
					Message: github.Ptr("Unexpected commit message"), // ⚠️ Не тот текст
				},
			},
		}, nil, nil).AnyTimes()

	hub := hub.NewHub(asyncProducer, slog.Default(), mockGitHubClient, "test_topic")
	hub.AddTrack("https://github.com/test_owner/test_repo")

	consumer, err := sarama.NewConsumer(addresses, nil)
	if !assert.NoError(t, err) {
		t.Fatalf("❌ Failed to create a New Consumer: %v", err)
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition("test_topic", 0, sarama.OffsetNewest)
	if !assert.NoError(t, err) {
		t.Fatalf("Failed to consume partition: %v", err)
	}
	defer partitionConsumer.Close()

	message := <-partitionConsumer.Messages()

	var msg map[string]interface{}
	if err := json.Unmarshal(message.Value, &msg); err != nil {
		t.Fatalf("❌ Failed to consume partition: %v", err)
	}

	commitMsg, ok1 := msg["message"].(string)
	sha, ok2 := msg["sha"].(string)

	if !ok1 || !ok2 {
		t.Fatalf("❌ Invalid message format: %v", msg)
	}

	// ❗ Ожидаем, что данные НЕ ДОЛЖНЫ совпасть, иначе тест провалится
	if commitMsg == "Test commit message" && sha == "test_sha" {
		t.Fatalf("❌ Test should have failed, but Kafka returned expected data! Message: %+v", msg)
	}

	t.Log("✅ Test passed: Received unexpected message, as expected")
}

func TestHub_EmptyCommit(t *testing.T) {
	addresses := []string{"0.0.0.0:9093"}

	asyncProducer, err := sarama.NewAsyncProducer(addresses, nil)
	if !assert.NoError(t, err) {
		t.Fatalf("❌ Failed to create AsyncProducer: %v", err)
	}
	defer asyncProducer.Close()

	ctrl := gomock.NewController(t)
	mockGitHubClient := git.NewMockGitHubClient(ctrl)
	// Возвращаем пустой список коммитов
	mockGitHubClient.EXPECT().ListCommits(gomock.Any(), "test_owner", "test_repo", gomock.Any()).Return(
		[]*github.RepositoryCommit{}, nil, nil).AnyTimes()

	hub := hub.NewHub(asyncProducer, slog.Default(), mockGitHubClient, "test_topic")
	hub.AddTrack("https://github.com/test_owner/test_repo")

	consumer, err := sarama.NewConsumer(addresses, nil)
	if !assert.NoError(t, err) {
		t.Fatalf("❌ Failed to create Consumer: %v", err)
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition("test_topic", 0, sarama.OffsetNewest)
	if !assert.NoError(t, err) {
		t.Fatalf("❌ Failed to consume partition: %v", err)
	}
	defer partitionConsumer.Close()

	select {
	case msg := <-partitionConsumer.Messages():
		t.Fatalf("❌ Unexpected message received: %+v", msg)
	default:
		t.Log("✅ Test passed: No messages received as expected")
	}
}

func TestHub_MultipleCommits(t *testing.T) {
	addresses := []string{"0.0.0.0:9093"}

	asyncProducer, err := sarama.NewAsyncProducer(addresses, nil)
	if !assert.NoError(t, err) {
		t.Fatalf("❌ Failed to create AsyncProducer: %v", err)
	}
	defer asyncProducer.Close()

	ctrl := gomock.NewController(t)
	mockGitHubClient := git.NewMockGitHubClient(ctrl)
	mockGitHubClient.EXPECT().ListCommits(gomock.Any(), "test_owner", "test_repo", gomock.Any()).Return(
		[]*github.RepositoryCommit{
			{
				SHA: github.Ptr("sha_1"),
				Commit: &github.Commit{
					Message: github.Ptr("First commit"),
				},
			},
			{
				SHA: github.Ptr("sha_2"),
				Commit: &github.Commit{
					Message: github.Ptr("Second commit"),
				},
			},
		}, nil, nil).AnyTimes()

	hub := hub.NewHub(asyncProducer, slog.Default(), mockGitHubClient, "test_topic")

	hub.AddTrack("https://github.com/test_owner/test_repo")

	consumer, err := sarama.NewConsumer(addresses, nil)
	if !assert.NoError(t, err) {
		t.Fatalf("❌ Failed to create Consumer: %v", err)
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition("test_topic", 0, sarama.OffsetNewest)
	if !assert.NoError(t, err) {
		t.Fatalf("❌ Failed to consume partition: %v", err)
	}
	defer partitionConsumer.Close()

	expectedCommits := map[string]string{
		"sha_1": "First commit",
	}

	receivedCommits := map[string]string{}

	message := <-partitionConsumer.Messages()
	var msg map[string]interface{}
	if err := json.Unmarshal(message.Value, &msg); err != nil {
		t.Fatalf("❌ Failed to unmarshal message: %v", err)
	}

	sha, _ := msg["sha"].(string)
	commitMsg, _ := msg["message"].(string)

	receivedCommits[sha] = commitMsg

	assert.Equal(t, expectedCommits, receivedCommits, "❌ Received commits do not match expected")
	t.Log("✅ Test passed: Received all expected commits")
}

func TestHub_GitHubError(t *testing.T) {
	addresses := []string{"0.0.0.0:9093"}

	asyncProducer, err := sarama.NewAsyncProducer(addresses, nil)
	if !assert.NoError(t, err) {
		t.Fatalf("❌ Failed to create AsyncProducer: %v", err)
	}
	defer asyncProducer.Close()

	ctrl := gomock.NewController(t)
	mockGitHubClient := git.NewMockGitHubClient(ctrl)
	mockGitHubClient.EXPECT().ListCommits(gomock.Any(), "test_owner", "test_repo", gomock.Any()).Return(
		nil, nil, assert.AnError).AnyTimes() // Возвращаем ошибку

	hub := hub.NewHub(asyncProducer, slog.Default(), mockGitHubClient, "test_topic")
	hub.AddTrack("https://github.com/test_owner/test_repo")

	consumer, err := sarama.NewConsumer(addresses, nil)
	if !assert.NoError(t, err) {
		t.Fatalf("❌ Failed to create Consumer: %v", err)
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition("test_topic", 0, sarama.OffsetNewest)
	if !assert.NoError(t, err) {
		t.Fatalf("❌ Failed to consume partition: %v", err)
	}
	defer partitionConsumer.Close()

	select {
	case msg := <-partitionConsumer.Messages():
		t.Fatalf("❌ Unexpected message received despite GitHub error: %+v", msg)
	default:
		t.Log("✅ Test passed: No messages received as expected")
	}
}
