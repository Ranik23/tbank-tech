package hub

import (
	"log/slog"
	"strconv"

	"testing"
	"time"

	mockgithub "tbank/scrapper/internal/mocks/github"

	"github.com/golang/mock/gomock"
	"github.com/google/go-github/v69/github"
	"github.com/stretchr/testify/assert"
)

func TestHub_AddLink_SendsCommit(t *testing.T) {

	commitChan := make(chan CustomCommit)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mockgithub.NewMockGitHubClient(ctrl)

	url := "https://github.com/Ranik23/tbank-tech"

	mockClient.EXPECT().LatestCommit(gomock.Any(), "Ranik23", "tbank-tech", gomock.Any()).Return(
		&github.RepositoryCommit{
			SHA: github.Ptr("test_sha"),
			Commit: &github.Commit{
				Message: github.Ptr("Test commit message"),
			},
		}, nil, nil,
	).AnyTimes()

	hub := NewHub(mockClient, commitChan, slog.Default())

	hub.AddLink(url, 1)


	select {
	case commit := <-commitChan:
		assert.Equal(t, "test_sha", *commit.Commit.SHA)
		assert.Equal(t, "Test commit message", *commit.Commit.Commit.Message)
	case <-time.After(5 * time.Second):
		t.Fatalf("Expected commit, but got none")
	}
}


func TestHub_AddLinks_SendsCommits(t *testing.T) {

	commitChan := make(chan CustomCommit)

	ctrl := gomock.NewController(t)

	mockClient := mockgithub.NewMockGitHubClient(ctrl)

	url := "https://github.com/Ranik23/tbank-tech"

	mockClient.EXPECT().LatestCommit(gomock.Any(), "Ranik23", "tbank-tech", gomock.Any()).Return(
		&github.RepositoryCommit{
			SHA: github.Ptr("test_sha"),
			Commit: &github.Commit{
				Message: github.Ptr("Test commit message"),
			},
		}, nil, nil,
	).AnyTimes()

	url2 := "https://github.com/Ranik23/weather-api-swagger"

	mockClient.EXPECT().LatestCommit(gomock.Any(), "Ranik23", "weather-api-swagger", gomock.Any()).Return(
		&github.RepositoryCommit{
			SHA: github.Ptr("test_sha2"),
			Commit: &github.Commit{
				Message: github.Ptr("Test commit message2"),
			},
		}, nil, nil,
	).AnyTimes()

	hub := NewHub(mockClient, commitChan, slog.Default())

	hub.AddLink(url, 1)
	hub.AddLink(url2, 2)

	count := 0
	expectedCount := 2
	
	outerLoop:
	for {
		select {
		case repoCommit := <-commitChan:
			t.Logf("SHA - %s, Message - %s", *repoCommit.Commit.SHA, *repoCommit.Commit.Commit.Message)
			count++
			if count == expectedCount {
				return
			}
		case <-time.After(5 * time.Second):
			break outerLoop
		}
	}
	
	assert.Equal(t, count, expectedCount)
}


func TestHub_Stop(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mockgithub.NewMockGitHubClient(ctrl)
	commitChan := make(chan CustomCommit)
	logger := slog.Default()

	hub := NewHub(mockClient, commitChan, logger)

	hub.AddLink("https://github.com/Ranik23/tbank-tech", 1)

	hub.Stop()

	select {
	case _, ok := <-commitChan:
		assert.False(t, ok, "Expected commitChan to be closed after Stop")
	case <-time.After(1 * time.Second):
		t.Fatalf("Expected commitChan to be closed, but it's still open")
	}
}

func TestWrongURLScheme(t *testing.T) {

	hub := NewHub(nil, nil, slog.Default())

	hub.AddLink("wjhfwfwfhw", 2)

	 _, ok := hub.pairCancelFunc.Load(Pair{"wjhfwfwfhw", "2"})
	if ok {
		t.Fatalf("No Goroutine Should Have Started Due To The Wrong URL Scheme")
	}
}

func TestHub_AddLink_ErrorFetchingCommit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mockgithub.NewMockGitHubClient(ctrl)
	commitChan := make(chan CustomCommit)
	logger := slog.Default()

	hub := NewHub(mockClient, commitChan, logger)

	url := "https://github.com/Ranik23/tbank-tech"

	mockClient.EXPECT().LatestCommit(gomock.Any(), "Ranik23", "tbank-tech", gomock.Any()).Return(
		nil, nil, assert.AnError,
	).AnyTimes()

	hub.AddLink(url, 1)

	time.Sleep(2 * time.Second)

	select {
	case <-commitChan:
		t.Fatalf("Expected no commit, but got one")
	default:

	}
}

func TestHub_AddLink_DuplicateLink(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mockgithub.NewMockGitHubClient(ctrl)
	commitChan := make(chan CustomCommit)
	logger := slog.Default()

	hub := NewHub(mockClient, commitChan, logger)

	url := "https://github.com/Ranik23/tbank-tech"
	userID := uint(1)

	mockClient.EXPECT().LatestCommit(gomock.Any(), "Ranik23", "tbank-tech", gomock.Any()).Return(
		&github.RepositoryCommit{
			SHA: github.Ptr("test_sha"),
			Commit: &github.Commit{
				Message: github.Ptr("Test commit message"),
			},
		}, nil, nil,
	).AnyTimes()

	
	hub.AddLink(url, userID)

	hub.AddLink(url, userID)

	pair := Pair{url, strconv.Itoa(int(userID))}
	_, ok := hub.pairCancelFunc.Load(pair)
	assert.True(t, ok, "Expected pair to be added to pairCancelFunc")

	select {
	case commit := <-commitChan:
		assert.Equal(t, "test_sha", *commit.Commit.SHA)
	case <-time.After(5 * time.Second):
		t.Fatalf("Expected commit, but got none")
	}

	// Проверяем, что в канале больше нет коммитов
	select {
	case <-commitChan:
		t.Fatalf("Expected no more commits, but got one")
	default:
		// Все в порядке, канал пуст
	}
}

func TestHub_RemoveLink_CancelsGoroutine(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mockgithub.NewMockGitHubClient(ctrl)
	commitChan := make(chan CustomCommit)
	logger := slog.Default()

	hub := NewHub(mockClient, commitChan, logger)

	url := "https://github.com/Ranik23/tbank-tech"
	userID := uint(1)

	mockClient.EXPECT().LatestCommit(gomock.Any(), "Ranik23", "tbank-tech", gomock.Any()).Return(
		&github.RepositoryCommit{
			SHA: github.Ptr("test_sha"),
			Commit: &github.Commit{
				Message: github.Ptr("Test commit message"),
			},
		}, nil, nil,
	).AnyTimes()

	hub.AddLink(url, userID)
	hub.RemoveLink(url, userID)

	time.Sleep(5 * time.Second)

	pair := Pair{url, strconv.Itoa(int(userID))}
	_, ok := hub.pairCancelFunc.Load(pair)
	assert.False(t, ok, "Expected pair to be removed from pairCancelFunc")
}

func TestHub_AddLink_UpdatesCommitOnNewSHA(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mockgithub.NewMockGitHubClient(ctrl)
	commitChan := make(chan CustomCommit)
	logger := slog.Default()

	hub := NewHub(mockClient, commitChan, logger)

	url := "https://github.com/Ranik23/tbank-tech"

	mockClient.EXPECT().LatestCommit(gomock.Any(), "Ranik23", "tbank-tech", gomock.Any()).Return(
		&github.RepositoryCommit{
			SHA: github.Ptr("test_sha1"),
			Commit: &github.Commit{
				Message: github.Ptr("Test commit message1"),
			},
		}, nil, nil,
	).Times(1)


	mockClient.EXPECT().LatestCommit(gomock.Any(), "Ranik23", "tbank-tech", gomock.Any()).Return(
		&github.RepositoryCommit{
			SHA: github.Ptr("test_sha2"),
			Commit: &github.Commit{
				Message: github.Ptr("Test commit message2"),
			},
		}, nil, nil,
	).Times(1)

	hub.AddLink(url, 1)


	select {
	case commit := <-commitChan:
		assert.Equal(t, "test_sha1", *commit.Commit.SHA)
		assert.Equal(t, "Test commit message1", *commit.Commit.Commit.Message)
	case <-time.After(5 * time.Second):
		t.Fatalf("Expected first commit, but got none")
	}

	select {
	case commit := <-commitChan:
		assert.Equal(t, "test_sha2", *commit.Commit.SHA)
		assert.Equal(t, "Test commit message2", *commit.Commit.Commit.Message)
	case <-time.After(5 * time.Second):
		t.Fatalf("Expected second commit, but got none")
	}
}

func TestHub_SyncMap_Operations(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mockgithub.NewMockGitHubClient(ctrl)
	commitChan := make(chan CustomCommit)
	logger := slog.Default()

	hub := NewHub(mockClient, commitChan, logger)

	url := "https://github.com/Ranik23/tbank-tech"
	userID := uint(1)


	hub.AddLink(url, userID)

	pair := Pair{url, strconv.Itoa(int(userID))}
	_, ok := hub.pairCancelFunc.Load(pair)
	assert.True(t, ok, "Expected pair to be added to pairCancelFunc")

	hub.RemoveLink(url, userID)

	_, ok = hub.pairCancelFunc.Load(pair)
	assert.False(t, ok, "Expected pair to be removed from pairCancelFunc")
}