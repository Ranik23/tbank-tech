package hub

import (
	"log/slog"
	"testing"
	"time"

	mockgithub "github.com/Ranik23/tbank-tech/scrapper/pkg/github_client/mock"

	"github.com/golang/mock/gomock"
	"github.com/google/go-github/v69/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHub_AddLink_UpdatesCommitOnNewSHA(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mockgithub.NewMockGitHubClient(ctrl)
	commitChan := make(chan CustomCommit)
	logger := slog.Default()

	hub := NewHub(mockClient, commitChan, logger)

	url := "https://github.com/Ranik23/tbank-tech"

	mockClient.EXPECT().LatestCommit(gomock.Any(), "Ranik23", "tbank-tech", gomock.Any(), gomock.Any()).Return(
		&github.RepositoryCommit{
			SHA: github.Ptr("test_sha1"),
			Commit: &github.Commit{
				Message: github.Ptr("Test commit message1"),
			},
		}, nil, nil,
	).Times(1)

	mockClient.EXPECT().LatestCommit(gomock.Any(), "Ranik23", "tbank-tech", gomock.Any(), gomock.Any()).Return(
		&github.RepositoryCommit{
			SHA: github.Ptr("test_sha2"),
			Commit: &github.Commit{
				Message: github.Ptr("Test commit message2"),
			},
		}, nil, nil,
	).Times(1)

	hub.Run()

	err := hub.AddLink(url, 1, "test_token", 4 * time.Second)
	require.NoError(t, err)


	select {
	case commit := <-commitChan:
		assert.Equal(t, "test_sha2", *commit.Commit.SHA)
		assert.Equal(t, "Test commit message2", *commit.Commit.Commit.Message)
	case <-time.After(5 * time.Second):
		t.Fatalf("Expected second commit, but got none")
	}

	hub.Stop()
}

func TestHub_AddLinks_SendsCommitsMultiple(t *testing.T) {

	commitChan := make(chan CustomCommit)

	ctrl := gomock.NewController(t)

	mockClient := mockgithub.NewMockGitHubClient(ctrl)

	url := "https://github.com/Ranik23/tbank-tech"

	mockClient.EXPECT().LatestCommit(gomock.Any(), "Ranik23", "tbank-tech", gomock.Any(), gomock.Any()).Return(
		&github.RepositoryCommit{
			SHA: github.Ptr("test_sha"),
			Commit: &github.Commit{
				Message: github.Ptr("Test commit message"),
			},
		}, nil, nil,
	).AnyTimes()

	url2 := "https://github.com/Ranik23/weather-api-swagger"

	mockClient.EXPECT().LatestCommit(gomock.Any(), "Ranik23", "weather-api-swagger", gomock.Any(), gomock.Any()).Return(
		&github.RepositoryCommit{
			SHA: github.Ptr("test_sha2"),
			Commit: &github.Commit{
				Message: github.Ptr("Test commit message2"),
			},
		}, nil, nil,
	).AnyTimes()

	hub := NewHub(mockClient, commitChan, slog.Default())

	hub.Run()

	hub.AddLink(url, 1, "test_token", 4*time.Second)
	hub.AddLink(url2, 2, "test_token", 4*time.Second)

	count := 0
	expectedCount := 0

outerLoop:
	for {
		select {
		case repoCommit := <-commitChan:
			t.Logf("SHA - %s, Message - %s", *repoCommit.Commit.SHA, *repoCommit.Commit.Commit.Message)
			count++
			if count == expectedCount {
				break outerLoop
			}
		case <-time.After(10 * time.Second):
			break outerLoop
		}
	}

	assert.Equal(t, count, expectedCount)

	hub.Stop()
}

func TestHub_AddLink_SendsCommit(t *testing.T) {

	commitChan := make(chan CustomCommit)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mockgithub.NewMockGitHubClient(ctrl)

	url := "https://github.com/Ranik23/tbank-tech"

	mockClient.EXPECT().LatestCommit(gomock.Any(), "Ranik23", "tbank-tech", gomock.Any(), gomock.Any()).Return(
		&github.RepositoryCommit{
			SHA: github.Ptr("test_sha"),
			Commit: &github.Commit{
				Message: github.Ptr("Test commit message"),
			},
		}, nil, nil,
	).Times(1)

	mockClient.EXPECT().LatestCommit(gomock.Any(), "Ranik23", "tbank-tech", gomock.Any(), gomock.Any()).Return(
		&github.RepositoryCommit{
			SHA: github.Ptr("test_sha2"),
			Commit: &github.Commit{
				Message: github.Ptr("Test commit message2"),
			},
		}, nil, nil,
	).AnyTimes()


	hub := NewHub(mockClient, commitChan, slog.Default())

	hub.Run()

	hub.AddLink(url, 1, "test_token", 10 * time.Second)

	select {
	case commit := <-commitChan:
		assert.Equal(t, "test_sha2", *commit.Commit.SHA)
		assert.Equal(t, "Test commit message2", *commit.Commit.Commit.Message)
	case <-time.After(20 * time.Second):
		t.Fatalf("Expected commit, but got none")
	}

	hub.Stop()
}

func TestHub_AddLink_ErrorFetchingCommit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mockgithub.NewMockGitHubClient(ctrl)
	commitChan := make(chan CustomCommit)
	logger := slog.Default()

	hub := NewHub(mockClient, commitChan, logger)

	url := "https://github.com/Ranik23/tbank-tech"

	mockClient.EXPECT().LatestCommit(gomock.Any(), "Ranik23", "tbank-tech", gomock.Any(), gomock.Any()).Return(
		nil, nil, assert.AnError,
	).AnyTimes()

	hub.Run()

	err := hub.AddLink(url, 1, "test_token", 10*time.Second)

	require.Error(t, err)


	hub.Stop()
}
