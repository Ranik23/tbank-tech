package hub

import (
	"encoding/json"
	"log/slog"
	gitmock "tbank/scrapper/internal/mocks/github"
	kafkamock "tbank/scrapper/internal/mocks/kafka"
	"testing"

	"github.com/IBM/sarama"
	"github.com/golang/mock/gomock"
	"github.com/google/go-github/v69/github"
	"github.com/stretchr/testify/assert"
)

func TestGetLatestCommit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGitHubClient := gitmock.NewMockGitHubClient(ctrl)
	mockGitHubClient.EXPECT().LatestCommit(gomock.Any(), "testOwner", "testRepo", gomock.Any()).
		Return(&github.RepositoryCommit{
				SHA: github.Ptr("test_sha"),
				Commit: &github.Commit{
					Message: github.Ptr("Test commit message"),
				},
			}, nil, nil).Times(1)

	linkToTrack := "https://github.com/testOwner/testRepo"

	client := NewClient(nil, slog.Default(), linkToTrack, mockGitHubClient)

	commit, err := client.getLatestCommit()

	assert.NoError(t, err)
	assert.NotNil(t, commit)
	assert.Equal(t, "test_sha", commit.GetSHA())
	assert.Equal(t, "Test commit message", commit.Commit.GetMessage())
}


func TestPublishCommit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAsyncProducer := kafkamock.NewMockAsyncProducer(ctrl)

	producerChan := make(chan *sarama.ProducerMessage, 1)
	mockAsyncProducer.EXPECT().Input().Return(producerChan).Times(1)


	linkToTrack := "https://github.com/testOwner/testRepo"
	
	client := NewClient(mockAsyncProducer, slog.Default(), linkToTrack, nil)

	commit := &github.RepositoryCommit{
		SHA: github.Ptr("test_sha"),
		Commit: &github.Commit{
			Message: github.Ptr("Test message"),
		},
	}

	client.publishCommit(commit)

	msg := <-producerChan
	assert.Equal(t, "testOwner/testRepo", msg.Topic)

	var payload map[string]string
	assert.NoError(t, json.Unmarshal(msg.Value.(sarama.ByteEncoder), &payload))
	assert.Equal(t, "test_sha", payload["sha"])
	assert.Equal(t, "Test message", payload["message"])
}
